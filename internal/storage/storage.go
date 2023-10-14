package storage

import (
	"fmt"
	"url-shortener/config"
)

type Storage interface {
	Save(user int, shortURL string, path string) error
	GetShortIfHave(user int, path string) (string, error)
	GetLong(user int, short string) (string, error)
	Close() error
	CheckConnection() error
	GetAll(user int) ([]URL, error)
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
