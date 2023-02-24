package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/pyropy/decentralised-lambda/executor/bacalhau"
	"github.com/pyropy/decentralised-lambda/job"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/system"
	shell "github.com/ipfs/go-ipfs-api"
)

var (
	DefaultTimeout = time.Second * 300
)

type Node struct {
	ipfsClient       *shell.Shell
	bacalhauExecutor *bacalhau.Executor
}

func NewNode(cfg *Config) (*Node, error) {
	ipfsClient := shell.NewShell("localhost:5001")

	err := system.InitConfig()
	if err != nil {
		return nil, err
	}

	bacalhauExecutor := bacalhau.NewExecutor(cfg.BacalhauEndpoint)

	return &Node{ipfsClient: ipfsClient, bacalhauExecutor: bacalhauExecutor}, nil
}

func (n *Node) InvokeLambdaFunction(ctx context.Context, wasmCid cid.Cid, request *http.Request) (interface{}, error) {
	cm := system.NewCleanupManager()
	inputCid, err := n.prepareJobInput(ctx, request)
	if err != nil {
		return nil, err
	}

	fmt.Println("inputCid", inputCid)
	j := job.NewJob("", wasmCid, inputCid)
	defer cm.Cleanup()

	resultCid, err := n.bacalhauExecutor.ExecuteJob(ctx, j)
	if err != nil {
		return nil, err
	}

	bodyCid := fmt.Sprintf("%s/outputs/output.json", resultCid.String())

	bodyBuf, err := n.ipfsClient.Cat(bodyCid)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(bodyBuf)
	if err != nil {
		return nil, err
	}

	var b map[string]interface{}

	err = json.Unmarshal(bodyBytes, &b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (n *Node) prepareJobInput(ctx context.Context, request *http.Request) (cid.Cid, error) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return cid.Cid{}, err
	}

	inputCidString, err := n.ipfsClient.Add(strings.NewReader(string(b)))
	if err != nil {
		return cid.Cid{}, err
	}

	inputCid, err := cid.Parse(inputCidString)
	if err != nil {
		return cid.Cid{}, err
	}

	return inputCid, nil
}
