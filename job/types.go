package job

import "github.com/ipfs/go-cid"

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
