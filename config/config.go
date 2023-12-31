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
	File    string
	DB      string
	Storage string
}

type F struct {
	host    *string
	baseURL *string
	file    *string
	db      *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=host")
	f.baseURL = flag.String("b", "http://localhost:8080/", "-b=base")
	f.file = flag.String("f", "qwe", "-f=storage")
	f.db = flag.String("d", "", "-d=db")
}

func New() (c Config) {
	flag.Parse()
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		f.host = &envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		f.baseURL = &envBaseAddr
	}
	if envFile := os.Getenv("FILE_STORAGE_PATH"); envFile != "" {
		f.file = &envFile
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	if *f.file != "qwe" {
		c.Storage = "file"
	}
	if *f.db != "" {
		c.Storage = "db"
	}
	c.Host = *f.host
	c.BaseURL = func() string {
		if strings.HasSuffix(*f.baseURL, "/") {
			return *f.baseURL
		}

		return fmt.Sprintf("%s/", *f.baseURL)
	}()
	c.File = *f.file
	c.DB = *f.db
	return c

}
