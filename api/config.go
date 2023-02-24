package api

type Config struct {
	Host string `envconfig:"API_HOST" default:"localhost"`
	Port string `envconfig:"API_PORT" default:"8080"`
}
