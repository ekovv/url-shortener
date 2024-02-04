package service

import (
	"fmt"
	"sync/atomic"
	"url-shortener/internal/storage"

	"github.com/google/uuid"
)

// SessionService sa
type SessionService struct {
	storage storage.Storage
	sessMap map[string]int
	count   *atomic.Uint64
}

// NewSessionService sa
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

// CreateIfNotExists sa
func (s *SessionService) CreateIfNotExists() (string, int) {
	s.count.Add(1)
	intID := int(s.count.Load())
	session := uuid.New().String()
	s.sessMap[session] = intID
	return session, intID
}

// GetID sa
func (s *SessionService) GetID(session string) int {
	id := s.sessMap[session]
	return id
}
