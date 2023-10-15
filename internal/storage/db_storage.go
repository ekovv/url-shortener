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

func (s *DBStorage) GetLastID() (int, error) {
	var lastID sql.NullInt64
	err := s.conn.QueryRow("SELECT MAX(id) FROM urls").Scan(&lastID)
	if err != nil {
		return 0, err
	}

	if lastID.Valid {
		return int(lastID.Int64), nil
	}

	return 0, nil
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

func (s *DBStorage) Save(user int, shortURL string, path string) error {
	insertQuery := `INSERT INTO urls(original, short, cookie) VALUES ($1, $2, $3)`
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

func (s *DBStorage) GetShortIfHave(user int, path string) (string, error) {
	query := "SELECT short FROM urls WHERE original = $1 AND cookie = $2"
	var short string
	err := s.conn.QueryRow(query, path, user).Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *DBStorage) GetLong(user int, short string) (string, error) {
	query := "SELECT original FROM urls WHERE short = $1"
	var original string
	err := s.conn.QueryRow(query, short).Scan(&original)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
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

func (s *DBStorage) GetAll(user int) ([]URL, error) {
	query := "SELECT original, short FROM urls WHERE cookie = $1"
	var list []URL
	rows, err := s.conn.Query(query, user)
	if err != nil {
		return nil, fmt.Errorf("failed to getall urls: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
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
