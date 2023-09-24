package config

import (
	"flag"
	"os"
)

type Config struct {
	Host    string
	BaseURL string
}

type F struct {
	host    *string
	baseURL *string
}

var f F

func init() {
	f.host = flag.String("a", "localhost:8080", "-a=host")
	f.baseURL = flag.String("b", "http://localhost:8080/", "-b=base")
}

func New() *Config {
	flag.Parse()
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		f.host = &envRunAddr
	}

	return &Config{
		Host:    *f.host,
		BaseURL: *f.baseURL,
	}
}
