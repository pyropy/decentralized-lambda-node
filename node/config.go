package node

type Config struct {
	IPFSEndpoint     string `envconfig:"IPFS_ENDPOINT"`
	BacalhauEndpoint string `envconfig:"BACALHAU_ENDPOINT"`
}
