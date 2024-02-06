package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"url-shortener/config"
	"url-shortener/internal/storage"

	_ "github.com/lib/pq"
	"github.com/speps/go-hashids/v2"
)

// Service struct
type Service struct {
	Storage storage.Storage
	config  config.Config
}

// NewService constructs a new Service
func NewService(storage storage.Storage, config config.Config) (Service, error) {
	return Service{
		Storage: storage,
		config:  config,
	}, nil
}

// GetShort save long
func (s *Service) GetShort(userID int, path string) (string, error) {
	short := s.getShortURL()
	err := s.Storage.Save(userID, short, path)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			short, err := s.Storage.GetShortIfHave(userID, path)
			if err != nil {
				return "", err
			}
			return s.config.BaseURL + short, storage.ErrAlreadyExists
		}
		return "", err
	}
	return s.config.BaseURL + short, nil
}

// GetLong get long with short
func (s *Service) GetLong(userID int, shortURL string) (string, error) {
	long, err := s.Storage.GetLong(userID, shortURL)
	if long == "" && err == nil {
		return "", nil
	}
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

// CheckConn check connection to db
func (s *Service) CheckConn() error {
	err := s.Storage.CheckConnection()
	if err != nil {
		return errors.New("not connected")
	}
	return nil
}

// SaveWithoutGenerate save links without generate
func (s *Service) SaveWithoutGenerate(userID int, id string, path string) (string, error) {
	err := s.Storage.Save(userID, id, path)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			short, err := s.Storage.GetShortIfHave(userID, path)
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
	hd.MinLength = 1
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

// GetAllUrls get all from db
func (s *Service) GetAllUrls(userID int) ([]storage.URL, error) {
	list, err := s.Storage.GetAll(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting %w", err)
	}
	return list, nil
}

// Delete delete urls in db
func (s *Service) Delete(list []string, id int) error {
	err := s.Storage.DeleteUrls(list, id)
	if err != nil {
		return fmt.Errorf("faile delete %w", err)
	}
	return nil
}
