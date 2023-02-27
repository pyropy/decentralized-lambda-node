package job

import "github.com/ipfs/go-cid"

var (
	ExecutionLayerBacalhau = "bacalhau"
)

var (
	PersistanceLayerIPFS    = "ipfs"
	PersistanceLayerEstuary = "estuary"
)

type JobSpec struct {
	// CID of the binary the binary to be executed.
	Binary           cid.Cid `json:"binary"`
	ExecutionLayer   string  `json:"executionLayer"`
	PersistanceLayer string  `json:"persistanceLayer"`
}

type Job struct {
	// ID of the job.
	ID string `json:"id"`
	// CID of the binary the binary to be executed.
	Binary cid.Cid `json:"binary"`
	// CID of the input data.
	Input cid.Cid `json:"input"`
}

func NewJob(id string, binary cid.Cid, input cid.Cid) *Job {
	return &Job{
		ID:     id,
		Binary: binary,
		Input:  input,
	}
}

func NewJobFromJobSpec(jobSpec *JobSpec) *Job {
	return &Job{
		Binary: jobSpec.Binary,
	}
}

func (j *Job) SetInput(input cid.Cid) {
	j.Input = input
}
