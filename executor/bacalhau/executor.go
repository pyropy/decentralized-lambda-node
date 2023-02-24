package bacalhau

import (
	"context"
	"fmt"
	bacalhauJob "github.com/filecoin-project/bacalhau/pkg/job"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/requester/publicapi"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/pyropy/decentralised-lambda/job"
	"log"
)

type Executor struct {
	bacalhauClient *publicapi.RequesterAPIClient
}

func NewExecutor(addr string) *Executor {
	bacalhauClient := publicapi.NewRequesterAPIClient(addr)
	return &Executor{
		bacalhauClient: bacalhauClient,
	}
}

func (e *Executor) ExecuteJob(ctx context.Context, job *job.Job) error {
	cm := system.NewCleanupManager()
	wasmJob := newWasmJob(job.Binary, job.Input)
	defer cm.Cleanup()
	j, err := ExecuteJob(ctx, e.bacalhauClient, cm, wasmJob)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("%+v", j))
	return nil
}

func ExecuteJob(
	ctx context.Context,
	apiClient *publicapi.RequesterAPIClient,
	cm *system.CleanupManager,
	j *model.Job) (*model.JobWithInfo, error) {

	err := bacalhauJob.VerifyJob(ctx, j)
	if err != nil {
		log.Fatal("Job failed to validate.")
		return nil, err
	}

	j, err = submitJob(ctx, apiClient, j)
	if err != nil {
		return nil, err
	}

	log.Println("Job ID: ", j.Metadata.ID)

	err = waitForJobToFinish(ctx, apiClient, j)
	if err != nil {
		return nil, err
	}

	jobReturn, found, err := apiClient.Get(ctx, j.Metadata.ID)
	if err != nil || !found {
		return nil, err
	}

	js, err := apiClient.GetJobState(ctx, jobReturn.Job.Metadata.ID)
	if err != nil {
		return nil, err
	}

	fmt.Println(js)

	// TODO: Download job results

	return jobReturn, nil
}
