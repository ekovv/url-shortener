package storage

import (
	"database/sql"
	"fmt"
	"url-shortener/config"
)

type DbStorage struct {
	conn *sql.DB
}

func (s *DbStorage) Save(shortURL string, path string) error {
	return nil
}

func (s *DbStorage) GetLong(short string) (string, error) {
	return "", nil
}

func (s *DbStorage) Close() error {
	return nil
}

func NewDBStorage(config config.Config) (*DbStorage, error) {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}
	defer db.Close()

	s := &DbStorage{}

	return s, s.CheckConnection()
}

func (s *DbStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}
