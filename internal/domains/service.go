package domains

import "url-shortener/internal/storage"

//go:generate go run github.com/vektra/mockery/v3 --name=UseCase
type UseCase interface {
	GetShort(user int, path string) (string, error)
	GetLong(user int, shortURL string) (string, error)
	CheckConn() error
	SaveWithoutGenerate(user int, id string, path string) (string, error)
	GetAllUrls(user int) ([]storage.URL, error)
	SaveAndGetSessionMap(session string) int
	Delete(list []string, id int) error
}
