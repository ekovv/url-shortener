package storage

type Storage struct {
	m map[string]string
}

func NewStorage() *Storage {
	return &Storage{m: make(map[string]string)}
}

func (s *Storage) Set() error {
	return nil
}
