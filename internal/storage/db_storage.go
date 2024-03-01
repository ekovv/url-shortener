package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

// DBStorage struct
type DBStorage struct {
	conn *sql.DB
}

// GetLastID get last id
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

// NewDBStorage constructor
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

// ErrAlreadyExists err
var ErrAlreadyExists = errors.New("already have")

// Save save in db
func (s *DBStorage) Save(user int, shortURL string, path string) error {
	insertQuery := `INSERT INTO urls(original, short, cookie, del) VALUES ($1, $2, $3, $4)`
	_, err := s.conn.Exec(insertQuery, path, shortURL, user, false)
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

// GetShortIfHave get short
func (s *DBStorage) GetShortIfHave(user int, path string) (string, error) {
	query := "SELECT short FROM urls WHERE original = $1 AND cookie = $2"
	var short string
	err := s.conn.QueryRow(query, path, user).Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

// GetLong get long
func (s *DBStorage) GetLong(userID int, short string) (string, error) {
	query := "SELECT original, del FROM urls WHERE short = $1"
	var original string
	var deleted bool
	err := s.conn.QueryRow(query, short).Scan(&original, &deleted)
	if deleted {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}
	return original, nil
}

// Close connection close
func (s *DBStorage) Close() error {
	return s.conn.Close()
}

// CheckConnection check connection
func (s *DBStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}

// GetAll get all
func (s *DBStorage) GetAll(userID int) ([]URL, error) {
	query := "SELECT original, short FROM urls WHERE cookie = $1"
	var list []URL
	rows, err := s.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to getall urls: %w", err)
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()
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

// DeleteUrls delete urls
func (s *DBStorage) DeleteUrls(list []string, userID int) error {
	for _, i := range list {
		query := "UPDATE urls SET del = true WHERE short = $1 AND cookie = $2"
		_, err := s.conn.Exec(query, i, userID)
		if err != nil {
			return fmt.Errorf("failed to delete %w", err)
		}
	}
	return nil
}
