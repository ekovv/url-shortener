package domains

//go:generate go run github.com/vektra/mockery/v3 --name=Usecase
type Usecase interface {
	RetShort(path string) (string, error)
	RetLong(shortURL string) (string, error)
}
