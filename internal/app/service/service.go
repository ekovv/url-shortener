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
	SetPair(path string, shortURL string) error
	GetLong(urlShort string) (string, error)
}

func (s *Service) GetShort(path string) (string, error) {
	short := s.getShortURL()
	err := s.storage.SetPair(path, short)
	if err != nil {
		return "", errors.New("invalid")
	}
	return s.config.BaseURL + short, nil
}

func (s *Service) GetLong(shortURL string) (string, error) {
	long, err := s.storage.GetLong(shortURL)
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

func (s *Service) getShortURL() string {
	randomString := generateRandomString(7)
	return randomString
}

func generateRandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}
