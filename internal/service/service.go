package service

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/speps/go-hashids/v2"
	"math/rand"
	"strconv"
	"time"
	"url-shortener/config"
	"url-shortener/internal/storage"
)

type Service struct {
	Storage storage.Storage
	config  config.Config
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
		if errors.Is(err, storage.ErrAlreadyExists) {
			short, err := s.Storage.GetShortIfHave(path)
			if err != nil {
				return "", err
			}
			return s.config.BaseURL + short, storage.ErrAlreadyExists
		}
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

func (s *Service) CheckConn() error {
	err := s.Storage.CheckConnection()
	if err != nil {
		return errors.New("not connected")
	}
	return nil
}

func (s *Service) SaveWithoutGenerate(id string, path string) (string, error) {
	err := s.Storage.Save(id, path)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			short, err := s.Storage.GetShortIfHave(path)
			if err != nil {
				return "", err
			}
			return s.config.BaseURL + short, storage.ErrAlreadyExists
		}
		return "", fmt.Errorf("failed to save in db: %w", err)
	}
	return s.config.BaseURL + id, nil
}

func (s *Service) getShortURL() string {
	hd := hashids.NewData()
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	var list []int
	for i := 0; i < 2; i++ {
		n := generator.Int63()
		s := strconv.FormatInt(n, 10)
		res, err := strconv.Atoi(s)
		if err != nil {
			return ""
		}
		list = append(list, res)
	}
	e, err := h.Encode(list)
	if err != nil {
		return ""
	}
	return e
}

//func generateRandomString(length int) string {
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//	randomInt := make([]byte, length)
//	for i := range randomInt {
//		randomInt[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(randomString)
//}
