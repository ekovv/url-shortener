package storage

import (
	"fmt"
	"url-shortener/config"
)

type Storage interface {
	Save(shortURL string, path string) error
	GetLong(short string) (string, error)
	Close() error
	CheckConnection() error
}

func New(cfg config.Config) (Storage, error) {
	switch cfg.Storage {
	case "db":
		d, err := NewDBStorage(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating db storage: %v", err)
		}
		return d, nil
	case "file":
		f, err := NewFileStorage(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("error creating db storage: %v", err)
		}
		return f, nil
	default:
		m := NewMapStorage()
		return m, nil
	}
}
