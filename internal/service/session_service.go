package service

import (
	"fmt"
	"sync/atomic"
	"url-shortener/internal/storage"
)

type SessionService struct {
	storage storage.Storage
	sessMap map[string]int
	count   *atomic.Uint64
}

func NewSessionService(storage storage.Storage) (SessionService, error) {
	lastID, err := storage.GetLastID()
	if err != nil {
		return SessionService{}, fmt.Errorf("error getting last id: %w", err)
	}

	newID := atomic.Uint64{}
	newID.Store(uint64(lastID))
	return SessionService{
		storage: storage,
		sessMap: make(map[string]int),
		count:   &newID,
	}, nil
}

func (s *SessionService) CreateIfNotExists(session string) int {
	a, ok := s.sessMap[session]
	if !ok {
		s.count.Add(1)
		intID := int(s.count.Load())
		s.sessMap[session] = intID
		return intID
	} else {
		return a
	}
}
