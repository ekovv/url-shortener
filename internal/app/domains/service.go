package domains

//go:generate go run github.com/vektra/mockery/v3 --name=UseCase
type UseCase interface {
	GetShort(path string) (string, error)
	GetLong(shortURL string) (string, error)
}
