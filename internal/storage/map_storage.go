package storage

import (
	"errors"
)

type MapStorage struct {
	m map[string]URL
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
		m: make(map[string]URL),
	}

}

func (s *MapStorage) Save(user string, shortURL string, path string) error {
	ur := URL{Original: path, Short: shortURL}
	s.m[user] = ur
	return nil
}

func (s *MapStorage) GetShortIfHave(user string, path string) (string, error) {
	ur, ok := s.m[user]
	if !ok {
		return "", errors.New("Short url not found")
	}

	return ur.Short, nil
}

func (s *MapStorage) GetLong(user string, urlShort string) (string, error) {
	ur, ok := s.m[user]
	if !ok {
		return "", errors.New("long url not found")
	}

	return ur.Original, nil
}

func (s *MapStorage) GetAll(user string) ([]URL, error) {
	var result []URL
	for _, value := range s.m {
		url := URL{}
		url.Original = value.Original
		url.Short = value.Short
		result = append(result, url)
	}
	return result, nil
}
