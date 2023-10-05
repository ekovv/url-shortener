package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
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
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate driver, %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"url", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to do migrate %w", err)
	}
	s := &DBStorage{
		conn: db,
	}

	return s, s.CheckConnection()
}

var ErrAlreadyExists = errors.New("already have")

func (s *DBStorage) Save(shortURL string, path string) error {
	insertQuery := `INSERT INTO urls(original, short) VALUES ($1, $2)`
	_, err := s.conn.Exec(insertQuery, path, shortURL)
	if err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code == "23505" {
				return ErrAlreadyExists
			}
			return err
		}
	}
	return nil
}

func (s *DBStorage) GetShortIfHave(path string) (string, error) {
	query := "SELECT short FROM urls WHERE original = $1"
	var short string
	err := s.conn.QueryRow(query, path).Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *DBStorage) GetLong(short string) (string, error) {
	query := "SELECT original FROM urls WHERE short = $1"
	var original string
	err := s.conn.QueryRow(query, short).Scan(&original)
	if err != nil {
		return "", err
	}
	return original, nil
}

func (s *DBStorage) Close() error {
	return s.conn.Close()
}

func (s *DBStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}
