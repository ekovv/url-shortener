package service

import (
	"errors"
	_ "github.com/lib/pq"
	"math/rand"
	"time"
	"url-shortener/config"
	"url-shortener/internal/app/storage"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Service struct {
	Storage     storage.Storage
	fileStorage storage.FileStorage
	config      config.Config
}

func NewService(storage storage.Storage, config config.Config) Service {
	return Service{
		Storage: storage,
		config:  config,
	}
}

func (s *Service) GetShort(path string) (string, error) {
	short := s.getShortURL()
	err := s.Storage.Save(short, path)
	if err != nil {
		return "", err
	}
	return s.config.BaseURL + short, nil
}

func (s *Service) GetLong(shortURL string) (string, error) {
	long, err := s.Storage.GetLong(shortURL)
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
