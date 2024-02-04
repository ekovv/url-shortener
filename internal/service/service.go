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

// Service sa
type Service struct {
	Storage storage.Storage
	config  config.Config
}

// NewService sa
func NewService(storage storage.Storage, config config.Config) (Service, error) {
	return Service{
		Storage: storage,
		config:  config,
	}, nil
}

// GetShort as
func (s *Service) GetShort(user int, path string) (string, error) {
	short := s.getShortURL()
	err := s.Storage.Save(user, short, path)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			short, err := s.Storage.GetShortIfHave(user, path)
			if err != nil {
				return "", err
			}
			return s.config.BaseURL + short, storage.ErrAlreadyExists
		}
		return "", err
	}
	return s.config.BaseURL + short, nil
}

// GetLong sa
func (s *Service) GetLong(user int, shortURL string) (string, error) {
	long, err := s.Storage.GetLong(user, shortURL)
	if long == "" && err == nil {
		return "", nil
	}
	if err != nil {
		return "", errors.New("invalid")
	}
	return long, nil
}

// CheckConn sa
func (s *Service) CheckConn() error {
	err := s.Storage.CheckConnection()
	if err != nil {
		return errors.New("not connected")
	}
	return nil
}

// SaveWithoutGenerate sa
func (s *Service) SaveWithoutGenerate(user int, id string, path string) (string, error) {
	err := s.Storage.Save(user, id, path)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			short, err := s.Storage.GetShortIfHave(user, path)
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

// GetAllUrls sa
func (s *Service) GetAllUrls(user int) ([]storage.URL, error) {
	list, err := s.Storage.GetAll(user)
	if err != nil {
		return nil, fmt.Errorf("error getting %w", err)
	}
	return list, nil
}

// Delete sa
func (s *Service) Delete(list []string, id int) error {
	err := s.Storage.DeleteUrls(list, id)
	if err != nil {
		return fmt.Errorf("faile delete %w", err)
	}
	return nil
}
