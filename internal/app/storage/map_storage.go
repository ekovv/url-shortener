package storage

import (
	"errors"
)

type Storage struct {
	m map[string]string
}

func NewStorage() *Storage {
	return &Storage{m: make(map[string]string)}
}

func (s *Storage) SetPair(path string, shortURL string) error {
	s.m[shortURL] = path
	return nil
}

func (s *Storage) GetLong(urlShort string) (string, error) {
	long, ok := s.m[urlShort]
	if !ok {
		return "", errors.New("long url not found")
	}

	return long, nil
}
