package storage

import (
	"errors"
)

type MapStorage struct {
	m map[string]string
}

func (s *MapStorage) Close() error {
	s.m = nil
	return nil
}

func NewMapStorage() *MapStorage {
	return &MapStorage{m: make(map[string]string)}
}

func (s *MapStorage) Save(shortURL string, path string) error {
	s.m[shortURL] = path
	return nil
}

func (s *MapStorage) GetLong(urlShort string) (string, error) {
	long, ok := s.m[urlShort]
	if !ok {
		return "", errors.New("long url not found")
	}

	return long, nil
}
