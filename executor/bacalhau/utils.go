package bacalhau

import (
	"context"
	"fmt"
	"github.com/filecoin-project/bacalhau/pkg/bacerrors"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/requester/publicapi"
	"github.com/ipfs/go-cid"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"time"
)

var (
	DefaultTimeout              = time.Second * 300
	HowFrequentlyToUpdateTicker = 50 * time.Millisecond
	AutoDownloadFolderPerm      = 0755
)

type eventStruct struct {
	Message    string
	IsTerminal bool
	IsError    bool
}

var terminalEvents = map[model.ExecutionStateType]eventStruct{
	// Need to add a carriage return to the end of the line, but only this one
	model.ExecutionStateFailed: {Message: "Error while executing the job.", IsTerminal: true, IsError: true},

	// Job is on StorageNode
	model.ExecutionStateResultRejected: {Message: "Results failed verification.", IsTerminal: true, IsError: false},
	model.ExecutionStateCompleted:      {Message: "", IsTerminal: true, IsError: false},

	// Job is canceled by the user
	model.ExecutionStateCanceled: {Message: "Job canceled by the user.", IsTerminal: true, IsError: true},
}

func newWasmJob(wasmCid cid.Cid, inputCid cid.Cid) *model.Job {
	wasmJob, _ := model.NewJobWithSaneProductionDefaults()
	wasmJob.Spec.Publisher = model.PublisherIpfs
	wasmJob.Spec.Engine = model.EngineWasm
	wasmJob.Spec.Verifier = model.VerifierNoop
	wasmJob.Spec.Timeout = DefaultTimeout.Seconds()
	wasmJob.Spec.Wasm.EntryPoint = "_start"
	wasmJob.Spec.Wasm.EnvironmentVariables = map[string]string{}
	wasmJob.Spec.Inputs = []model.StorageSpec{
		{
			StorageSource: model.StorageSourceIPFS,
			CID:           inputCid.String(),
			Path:          "/inputs/input.json",
		},
	}
	wasmJob.Spec.Outputs = []model.StorageSpec{
		{
			Name: "outputs",
			Path: "/outputs",
		},
	}
	wasmJob.Spec.Wasm.EntryModule = model.StorageSpec{
		StorageSource: model.StorageSourceIPFS,
		CID:           wasmCid.String(),
	}
	return wasmJob
}

func submitJob(
	ctx context.Context,
	apiClient *publicapi.RequesterAPIClient,
	j *model.Job,
) (*model.Job, error) {
	j, err := apiClient.Submit(ctx, j)
	if err != nil {
		return &model.Job{}, errors.Wrap(err, "failed to submit job")
	}
	return j, err
}

func waitForJobToFinish(ctx context.Context, apiClient *publicapi.RequesterAPIClient, j *model.Job) error {
	if j == nil || j.Metadata.ID == "" {
		return errors.New("No job returned from the server.")
	}

	time.Sleep(1 * time.Second)

	jobEvents, err := apiClient.GetEvents(ctx, j.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failure retrieving job events '%s': %s\n", j.Metadata.ID, err)
	}

	var ticker *time.Ticker
	var tickerDone = make(chan bool)
	var ShutdownSignals = []os.Signal{
		os.Interrupt,
	}

	ticker = time.NewTicker(HowFrequentlyToUpdateTicker)

	// Capture Ctrl+C if the user wants to finish early the job
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, ShutdownSignals...)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	finishedRunning := false
	var returnError error
	returnError = nil

	// goroutine for handling spinner ticks and spinner completion messages
	go func() {
		for {
			select {
			case <-tickerDone:
				ticker.Stop()
				return
			case <-ticker.C:
			}
		}
	}()

	// goroutine for handling SIGINT from the signal channel, or context
	// completion messages.
	go func() {
		select {
		case s := <-signalChan: // first signal, cancel context
			if s == os.Interrupt {
				if !finishedRunning {
					returnError = fmt.Errorf("Received signal. Exiting.")
				}
			} else {
				fmt.Println("Unexpected signal received. Exiting.")
			}
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	// Loop through the events, printing those that are interesting, and then
	// shutting down when a this job reaches a terminal state.
	for {

		if err != nil {
			if _, ok := err.(*bacerrors.JobNotFound); ok {
				returnError = fmt.Errorf(`Somehow even though we submitted a job successfully,
											we were not able to get its status. ID: %s`, j.Metadata.ID)
			} else {
				returnError = fmt.Errorf("Unknown error trying to get job (ID: %s): %+v", j.Metadata.ID, err)
			}

			finishedRunning = true

			tickerDone <- true
			signalChan <- os.Interrupt
			return returnError
		}

		for i := range jobEvents {
			if terminalEvents[jobEvents[i].NewStateType].IsTerminal {
				// Send a signal to the goroutine that is waiting for Ctrl+C
				finishedRunning = true

				tickerDone <- true
				signalChan <- os.Interrupt
				return err
			}
		}

		if condition := ctx.Err(); condition != nil {
			signalChan <- os.Interrupt
			break
		} else {
			jobEvents, err = apiClient.GetEvents(ctx, j.Metadata.ID)
			if err != nil {
				if _, ok := err.(*bacerrors.ContextCanceledError); ok {
					// We're done, the user canceled the job
					break
				} else {
					return errors.Wrap(err, "Error getting job events")
				}
			}
		}

		time.Sleep(time.Duration(500) * time.Millisecond) //nolint:gomnd // 500ms sleep
	} // end for

	return returnError
}
