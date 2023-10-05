package storage

import (
	"errors"
	"go.uber.org/zap"
)

type MapStorage struct {
	m      map[string]string
	logger zap.Logger
}

func (s *MapStorage) CheckConnection() error {
	if s.m == nil {
		return errors.New("map nil")
	}
	return nil
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

func (s *MapStorage) GetShortIfHave(path string) (string, error) {
	short, ok := s.m[path]
	if !ok {
		return "", errors.New("short url not found")
	}

	return short, nil
}

func (s *MapStorage) GetLong(urlShort string) (string, error) {
	long, ok := s.m[urlShort]
	if !ok {
		return "", errors.New("long url not found")
	}

	return long, nil
}
