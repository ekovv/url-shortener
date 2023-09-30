package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host    string
	BaseURL string
	Storage string
}

type F struct {
	host    *string
	baseURL *string
	storage *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=host")
	f.baseURL = flag.String("b", "http://localhost:8080/", "-b=base") // TODO: add const
	f.storage = flag.String("f", "map", "-f=storage")                 // TODO: add const
}

func New() *Config {
	flag.Parse()
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		f.host = &envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		f.baseURL = &envBaseAddr
	}
	if envStorage := os.Getenv("FILE_STORAGE_PATH"); envStorage != "" {
		f.storage = &envStorage
	}

	return &Config{
		Host: *f.host,
		BaseURL: func() string {
			if strings.HasSuffix(*f.baseURL, "/") {
				return *f.baseURL
			}

			return fmt.Sprintf("%s/", *f.baseURL)
		}(),
		Storage: *f.storage,
	}
}
