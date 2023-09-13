package storage

import "errors"

type Storage struct {
	m map[string]string
}

func NewStorage() *Storage {
	return &Storage{m: make(map[string]string)}
}

func (s *Storage) Set(path string, shortUrl string) error {
	s.m[path] = shortUrl
	return nil
}

func (s *Storage) GetShort(path string) (string, error) {
	str, ok := s.m[path]
	if !ok {
		return "", errors.New("invalid original url")
	}
	return str, nil
}

func (s *Storage) GetLong(urlShort string) (string, error) {
	for key, val := range s.m {
		if val == urlShort {
			return key, nil
		}
	}
	return "", errors.New("invalid")
}
