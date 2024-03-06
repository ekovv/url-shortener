package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Config sa
type Config struct {
	Host    string `json:"server_address"`
	BaseURL string `json:"base_url"`
	TLS     bool   `json:"enable_https"`
	File    string `json:"file_storage_path"`
	DB      string `json:"database_dsn"`
	Storage string
	CFile   string
}

// F sa
type F struct {
	host    *string
	baseURL *string
	tls     *bool
	file    *string
	db      *string
	cFile   *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=host")
	f.baseURL = flag.String("b", "http://localhost:8080/", "-b=base")
	f.tls = flag.Bool("s", false, "-t=")
	f.file = flag.String("f", "qwe", "-f=storage")
	f.db = flag.String("d", "", "-d=db")
	f.cFile = flag.String("c", "", "-c=")
}

// New sa
func New() (c Config) {
	flag.Parse()
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		f.host = &envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		f.baseURL = &envBaseAddr
	}
	if envTLS := os.Getenv("ENABLE_HTTPS"); envTLS != "" {
		f.baseURL = &envTLS
	}
	if envFile := os.Getenv("FILE_STORAGE_PATH"); envFile != "" {
		f.file = &envFile
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		f.cFile = &envConfig
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
	c.TLS = *f.tls
	c.DB = *f.db
	c.CFile = *f.cFile
	file, err := os.Open(c.CFile)
	if err != nil {
		return
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(all, &c)
	if err != nil {
		return
	}
	return c

}
