package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/pyropy/decentralised-lambda/job"
	ufcli "github.com/urfave/cli/v2"
)

var WasmCmd = ufcli.Command{
	Name: "wasm",
	Subcommands: []*ufcli.Command{
		&WasmDeployCmd,
	},
}

var WasmDeployCmd = ufcli.Command{
	Name:  "deploy",
	Usage: "Upload WASM executable to IPFS",
	Flags: []ufcli.Flag{
		&ufcli.StringFlag{
			Name:  "file",
			Usage: "path to WASM file",
		},
		&ufcli.StringFlag{
			Name:  "ipfs",
			Usage: "IPFS Multiaddr API endpoint",
			Value: "/ip4/0.0.0.0/tcp/5001",
		},
	},
	Action: func(cctx *ufcli.Context) error {
		file := cctx.String("file")
		ipfsClient := shell.NewShell("localhost:5001")

		x, err := ipfsClient.AddLink(file)
		if err != nil {
			return err
		}

		spec := job.JobSpec{
			Binary:           cid.MustParse(x),
			ExecutionLayer:   job.ExecutionLayerBacalhau,
			PersistanceLayer: job.PersistanceLayerIPFS,
		}

		b, err := json.Marshal(spec)
		if err != nil {
			return err
		}

		specCid, err := ipfsClient.Add(strings.NewReader(string(b)))
		if err != nil {
			return err
		}

		fmt.Println(specCid)

		return nil
	},
}
