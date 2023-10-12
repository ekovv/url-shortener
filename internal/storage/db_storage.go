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

func (s *DBStorage) Save(user string, shortURL string, path string) error {
	insertQuery := `INSERT INTO urls(Original, Short, cookie) VALUES ($1, $2, $3)`
	_, err := s.conn.Exec(insertQuery, path, shortURL, user)
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

func (s *DBStorage) GetShortIfHave(user string, path string) (string, error) {
	query := "SELECT Short FROM urls WHERE Original = $1"
	var short string
	err := s.conn.QueryRow(query, path, user).Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *DBStorage) GetLong(user string, short string) (string, error) {
	query := "SELECT Original FROM urls WHERE Short = $1 AND cookie = $2"
	var original string
	err := s.conn.QueryRow(query, short, user).Scan(&original)
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

func (s *DBStorage) GetAll(user string) ([]URL, error) {
	query := "SELECT Original, Short FROM urls WHERE cookie = $1"
	var list []URL
	rows, err := s.conn.Query(query, user)
	if err != nil {
		return nil, fmt.Errorf("failed to getall urls: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		url := URL{}
		err = rows.Scan(&url.Original, &url.Short)
		if err != nil {
			return nil, fmt.Errorf("failed to getall urls: %w", err)
		}
		list = append(list, url)
	}
	return list, nil
}
