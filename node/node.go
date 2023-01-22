package node

import (
	"context"
	"encoding/json"
	"github.com/filecoin-project/bacalhau/pkg/job"
	"github.com/filecoin-project/bacalhau/pkg/requester/publicapi"
	"github.com/ipfs/go-cid"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/ipfs"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/system"
)

var DefaultTimeout = time.Second * 300

type Node struct {
	ipfsClient     *ipfs.Client
	bacalhauClient *publicapi.RequesterAPIClient
}

func NewNode(cfg *Config) (*Node, error) {
	ipfsClient, err := ipfs.NewClient(cfg.IPFSEndpoint)
	if err != nil {
		return nil, err
	}

	err = system.InitConfig()
	if err != nil {
		return nil, err
	}

	bacalhauClient := publicapi.NewRequesterAPIClient(cfg.BacalhauEndpoint)

	return &Node{ipfsClient: ipfsClient, bacalhauClient: bacalhauClient}, nil
}

func newWasmJob(wasmCid cid.Cid, inputCid cid.Cid) *model.Job {
	wasmJob, _ := model.NewJobWithSaneProductionDefaults()
	wasmJob.Spec.Engine = model.EngineWasm
	wasmJob.Spec.Verifier = model.VerifierNoop
	wasmJob.Spec.Timeout = DefaultTimeout.Seconds()
	wasmJob.Spec.Wasm.EntryPoint = "_start"
	wasmJob.Spec.Wasm.EnvironmentVariables = map[string]string{}
	wasmJob.Spec.Publisher = model.PublisherIpfs
	wasmJob.Spec.Outputs = []model.StorageSpec{
		{
			Name: "outputs",
			Path: "/outputs",
		},
	}
	wasmJob.Spec.Contexts = append(wasmJob.Spec.Contexts, model.StorageSpec{
		StorageSource: model.StorageSourceIPFS,
		CID:           wasmCid.String(),
		Path:          "/job",
	})
	return wasmJob
}

func (n *Node) InvokeLambdaFunction(ctx context.Context, wasmCID cid.Cid, request *http.Request) error {
	cm := system.NewCleanupManager()
	inputCid, err := n.prepareJobInput(ctx, request)
	if err != nil {
		return err
	}

	wasmJob := newWasmJob(wasmCID, inputCid)
	defer cm.Cleanup()

	j, err := ExecuteJob(ctx, n.bacalhauClient, cm, wasmJob)
	if err != nil {
		return err
	}

	log.Println(j)
	return nil
}

func (n *Node) prepareJobInput(ctx context.Context, request *http.Request) (cid.Cid, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return cid.Cid{}, nil
	}

	tmpFile, err := os.CreateTemp("", "input*.json")
	if err != nil {
		return cid.Cid{}, nil
	}

	defer func() {
		err := os.RemoveAll(tmpFile.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = tmpFile.Write(b)
	if err != nil {
		return cid.Cid{}, nil
	}

	inputCidString, err := n.ipfsClient.Put(ctx, tmpFile.Name())
	if err != nil {
		return cid.Cid{}, nil
	}

	inputCid, err := cid.Parse(inputCidString)
	if err != nil {
		return cid.Cid{}, err
	}

	return inputCid, nil
}

func ExecuteJob(
	ctx context.Context,
	apiClient *publicapi.RequesterAPIClient,
	cm *system.CleanupManager,
	j *model.Job) (*model.Job, error) {

	err := job.VerifyJob(ctx, j)
	if err != nil {
		log.Fatal("Job failed to validate.")
		return nil, err
	}

	j, err = submitJob(ctx, apiClient, j)
	if err != nil {
		return nil, err
	}

	jobReturn, found, err := apiClient.Get(ctx, j.Metadata.ID)
	if err != nil || !found {
		return nil, err
	}

	_, err = apiClient.GetJobState(ctx, jobReturn.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return jobReturn, nil
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

//func downloadJobResult(
//	ctx context.Context,
//	cm *system.CleanupManager,
//	jobID string,
//	downloadSettings ipfs.IPFSDownloadSettings,
//)
