package executor

import (
	"context"
	"github.com/ipfs/go-cid"
	"github.com/pyropy/decentralised-lambda/job"
)

type Executor interface {
	// ExecuteJob executes a job and returns result Cid or error.
	ExecuteJob(ctx context.Context, j job.Job) (cid.Cid, error)
}
