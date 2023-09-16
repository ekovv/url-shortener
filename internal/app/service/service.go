package service

import (
	"errors"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return Service{storage: storage}
}

type Storage interface {
	Set(path string, shortURL string) error
	GetShort(path string) (string, error)
	GetLong(urlShort string) (string, error)
}

func (s *Service) RetShort(path string) (string, error) {
	err := s.storage.Set(path, ReplaceURLOnShort())
	if err != nil {
		return "", errors.New("invalid")
	}
	short, err := s.storage.GetShort(path)
	if err != nil {
		return "", errors.New("invalid")
	}
	return short, nil
}

func (s *Service) RetLong(shortURL string) (string, error) {
	long, err := s.storage.GetLong(shortURL)
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

func ReplaceURLOnShort() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortLink := shortenLink()
	shortLink = "http://localhost:8080/" + shortLink
	return shortLink
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
