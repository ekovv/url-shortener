package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/speps/go-hashids/v2"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
	"url-shortener/config"
	"url-shortener/internal/storage"
)

type Service struct {
	Storage storage.Storage
	config  config.Config
	sessMap map[string]int
	count   *atomic.Uint64
}

func NewService(storage storage.Storage, config config.Config) (Service, error) {
	lastID, err := storage.GetLastID()
	if err != nil {
		return Service{}, fmt.Errorf("error getting last id: %w", err)
	}

	newID := atomic.Uint64{}
	newID.Store(uint64(lastID))
	return Service{
		Storage: storage,
		config:  config,
		sessMap: make(map[string]int),
		count:   &newID,
	}, nil
}

func GenerateUUID() string {
	newToken := uuid.New().String()
	return newToken
}

func (s *Service) SaveAndGetSessionMap(session string) int {
	a, ok := s.sessMap[session]
	if !ok {
		s.count.Add(1)
		intID := int(s.count.Load())
		s.sessMap[session] = intID
		return intID
	} else {
		return a
	}
}

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

func (s *Service) GetLong(user int, shortURL string) (string, error) {
	long, err := s.Storage.GetLong(user, shortURL)
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

func (s *Service) GetAllUrls(user int) ([]storage.URL, error) {
	list, err := s.Storage.GetAll(user)
	if err != nil {
		return nil, fmt.Errorf("error getting %w", err)
	}
	return list, nil
}
