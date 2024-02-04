package storage

import (
	"errors"
	"fmt"
)

type MapStorage struct {
	m map[string]URLInfo
}

// DeleteUrls sa
func (s *MapStorage) DeleteUrls(list []string, user int) error {
	//TODO implement me
	panic("implement me")
}

// GetLastID sa
func (s *MapStorage) GetLastID() (int, error) {
	return len(s.m), nil
}

// CheckConnection sa
func (s *MapStorage) CheckConnection() error {
	if s.m == nil {
		return errors.New("map nil")
	}
	return nil
}

// Close sa
func (s *MapStorage) Close() error {
	s.m = nil
	return nil
}

// NewMapStorage sa
func NewMapStorage() *MapStorage {
	return &MapStorage{
		m: make(map[string]URLInfo),
	}

}

// Save sa
func (s *MapStorage) Save(user int, shortURL string, path string) error {
	ur := URLInfo{Original: path, User: user}
	s.m[shortURL] = ur
	return nil
}

// GetShortIfHave sa
func (s *MapStorage) GetShortIfHave(user int, path string) (string, error) {
	for key, value := range s.m {
		if value.Original == path && value.User == user {
			return key, nil
		}
	}
	return "", fmt.Errorf("not found")
}

// GetLong sa
func (s *MapStorage) GetLong(_ int, urlShort string) (string, error) {
	ur, ok := s.m[urlShort]
	if !ok {
		return "", errors.New("long url not found")
	}

	return ur.Original, nil
}

// GetAll sa
func (s *MapStorage) GetAll(user int) ([]URL, error) {
	var result []URL
	for key, value := range s.m {
		if user != value.User {
			continue
		}
		url := URL{}
		url.Original = value.Original
		url.Short = key
		result = append(result, url)
	}
	return result, nil
}
