package storage

import (
	"errors"
	"fmt"
)

type MapStorage struct {
	m map[string]URLInfo
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
	return &MapStorage{
		m: make(map[string]URLInfo),
	}

}

func (s *MapStorage) Save(user string, shortURL string, path string) error {
	ur := URLInfo{Original: path, User: user}
	s.m[shortURL] = ur
	return nil
}

func (s *MapStorage) GetShortIfHave(user string, path string) (string, error) {
	for key, value := range s.m {
		if value.Original == path {
			return key, nil
		}
	}
	return "", fmt.Errorf("Not found")
}

func (s *MapStorage) GetLong(_ string, urlShort string) (string, error) {
	ur, ok := s.m[urlShort]
	if !ok {
		return "", errors.New("long url not found")
	}

	return ur.Original, nil
}

func (s *MapStorage) GetAll(user string) ([]URL, error) {
	var result []URL
	for key, value := range s.m {
		url := URL{}
		url.Original = value.Original
		url.Short = key
		result = append(result, url)
	}
	return result, nil
}
