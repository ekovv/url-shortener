package storage

import (
	"database/sql"
	"fmt"
	"url-shortener/config"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(config config.Config) (*DBStorage, error) {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}
	defer db.Close()

	s := &DBStorage{}

	return s, s.CheckConnection()
}

func (s *DBStorage) Save(shortURL string, path string) error {
	return nil
}

func (s *DBStorage) GetLong(short string) (string, error) {
	return "", nil
}

func (s *DBStorage) Close() error {
	return nil
}

func (s *DBStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}
