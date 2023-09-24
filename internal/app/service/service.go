package service

import (
	"errors"
	"math/rand"
	"time"
	"url-shortener/config"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Service struct {
	storage Storage
	config  config.Config
}

func NewService(storage Storage, config config.Config) Service {
	return Service{
		storage: storage,
		config:  config,
	}
}

type Storage interface {
	Set(path string, shortURL string) error
	GetShort(path string) (string, error)
	GetLong(urlShort string) (string, error)
}

func (s *Service) RetShort(path string) (string, error) {
	short := s.ReplaceURLOnShort()
	err := s.storage.Set(path, short)
	if err != nil {
		return "", errors.New("invalid")
	}
	return s.config.BaseURL + short, nil
}

func (s *Service) RetLong(shortURL string) (string, error) {
	long, err := s.storage.GetLong(shortURL)
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

func (s *Service) ReplaceURLOnShort() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return shortenLink()
}

func generateRandomString(length int) string {
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

func shortenLink() string {
	randomString := generateRandomString(7)
	return randomString
}
