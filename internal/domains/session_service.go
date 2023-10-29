package domains

//go:generate go run github.com/vektra/mockery/v3 --name=SessionUseCase
type SessionUseCase interface {
	CreateIfNotExists(session string) int
}
