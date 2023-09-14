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
	Set(path string, shortUrl string) error
	GetShort(path string) (string, error)
	GetLong(urlShort string) (string, error)
}

func (s *Service) RetShort(path string) (string, error) {
	err := s.storage.Set(path, ReplaceUrlOnShort(path))
	if err != nil {
		return "", errors.New("invalid")
	}
	short, err := s.storage.GetShort(path)
	if err != nil {
		return "", errors.New("invalid")
	}
	return short, nil
}

func (s *Service) RetLong(shortUrl string) (string, error) {
	long, err := s.storage.GetLong(shortUrl)
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

func ReplaceUrlOnShort(path string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	url := path
	shortLink := shortenLink(url)
	return shortLink
}

func generateRandomString(length int) string {
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

func shortenLink(url string) string {
	shortLink := url
	if len(url) > 7 && url[:7] == "http://" {
		randomString := generateRandomString(7)
		shortLink = "http://" + randomString
	}

	return shortLink
}
