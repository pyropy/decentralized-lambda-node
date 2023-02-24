package node

type Config struct {
	IPFSEndpoint     string `envconfig:"IPFS_ENDPOINT"`
	BacalhauEndpoint string `envconfig:"BACALHAU_ENDPOINT"`
}

func DefaultConfig() *Config {
	return &Config{
		IPFSEndpoint:     "/ip4/0.0.0.0/tcp/5001",
		BacalhauEndpoint: "http://bootstrap.development.bacalhau.org:1234",
	}
}
