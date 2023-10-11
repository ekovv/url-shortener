package storage

import (
	"errors"
)

type MapStorage struct {
	m map[string]map[string]string
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
	return &MapStorage{m: make(map[string]map[string]string)}
}

func (s *MapStorage) Save(user string, shortURL string, path string) error {
	s.m[user][shortURL] = path
	return nil
}

func (s *MapStorage) GetShortIfHave(user string, path string) (string, error) {
	short, ok := s.m[user][path]
	if !ok {
		return "", errors.New("Short url not found")
	}

	return short, nil
}

func (s *MapStorage) GetLong(user string, urlShort string) (string, error) {
	long, ok := s.m[user][urlShort]
	if !ok {
		return "", errors.New("long url not found")
	}

	return long, nil
}

func (s *MapStorage) GetAll(user string) ([]URL, error) {
	var result []URL
	for key, value := range s.m[user] {
		url := URL{}
		url.Original = value
		url.Short = key
		result = append(result, url)
	}
	return result, nil
}
