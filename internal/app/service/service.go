package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"url-shortener/config"
	"url-shortener/internal/app/storage"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Service struct {
	mapStorage  MapStorage
	fileStorage storage.FileStorage
	config      config.Config
}

func NewService(storageMap MapStorage, storageFile storage.FileStorage, config config.Config) Service {
	return Service{
		mapStorage:  storageMap,
		fileStorage: storageFile,
		config:      config,
	}
}

type MapStorage interface {
	SetPair(path string, shortURL string) error
	GetLong(urlShort string) (string, error)
}

func (s *Service) GetShort(path string) (string, error) {
	short := s.getShortURL()
	if s.config.Storage != "map" {
		err := s.fileStorage.SaveInFile(short, path)
		if err != nil {
			fmt.Println("Not save in file")
			return "", err
		}
		return short, nil
	}
	err := s.mapStorage.SetPair(path, short)
	if err != nil {
		return "", errors.New("invalid")
	}
	return s.config.BaseURL + short, nil
}

func (s *Service) GetLong(shortURL string) (string, error) {
	long, err := s.mapStorage.GetLong(shortURL)
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
