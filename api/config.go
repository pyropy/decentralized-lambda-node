package api

type Config struct {
	Host string `envconfig:"API_HOST" default:"localhost"`
	Port string `envconfig:"API_PORT" default:"8080"`
}

func DefaultConfig() *Config {
	return &Config{
		Host: "localhost",
		Port: "6969",
	}
}
