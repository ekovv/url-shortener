package domains

import "url-shortener/internal/storage"

//go:generate go run github.com/vektra/mockery/v3 --name=UseCase
type UseCase interface {
	GetShort(user string, path string) (string, error)
	GetLong(user string, shortURL string) (string, error)
	CheckConn() error
	SaveWithoutGenerate(user string, id string, path string) (string, error)
	GetAllUrls(user string) ([]storage.URL, error)
}
