package service

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
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
		return s.config.BaseURL + short, nil
	}
	err := s.mapStorage.SetPair(path, short)
	if err != nil {
		return "", errors.New("invalid")
	}
	return s.config.BaseURL + short, nil
}

func (s *Service) GetLong(shortURL string) (string, error) {
	if s.config.Storage != "map" {
		long, err := s.fileStorage.GetLong(shortURL)
		if err != nil {
			fmt.Println("file without short")
			return "", err
		}
		return long, nil
	}
	long, err := s.mapStorage.GetLong(shortURL)
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

//func (s *Service) CheckConnection(conn string) error {
//	db, err := sql.Open("postgres", "postgres://username:password@localhost/dbname?sslmode=disable")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	err = db.Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//	return nil
//}

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
