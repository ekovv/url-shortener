package domains

//go:generate go run github.com/vektra/mockery/v3 --name=SessionUseCase
type SessionService interface {
	CreateIfNotExists() (string, int)
	GetID(session string) int
}
