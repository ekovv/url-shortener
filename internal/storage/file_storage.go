package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// FileStorage struct
type FileStorage struct {
	Path  string
	File  *os.File
	count int
}

// DeleteUrls delete
func (s *FileStorage) DeleteUrls(list []string, user int) error {
	//TODO implement me
	panic("implement me")
}

// GetLastID get last id
func (s *FileStorage) GetLastID() (int, error) {
	scanner := bufio.NewScanner(s.File)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count, nil
}

// CheckConnection check connection
func (s *FileStorage) CheckConnection() error {
	if s.File == nil {
		return fmt.Errorf("file not open")
	}
	return nil

}

// NewFileStorage constructor
func NewFileStorage(path string) (*FileStorage, error) {
	fs := &FileStorage{
		Path:  path,
		count: 1,
	}
	err := fs.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening file storage %w", err)
	}

	return fs, nil
}

// Open  file
func (s *FileStorage) Open() error {
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	s.File = file
	return nil
}

// ShutDown file
func (s *FileStorage) ShutDown() error {
	if err := s.File.Close(); err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}
	return nil
}

// Close fileCLose
func (s *FileStorage) Close() error {
	return s.File.Close()
}

// Save save in file
func (s *FileStorage) Save(user int, short string, long string) error {
	flag := false
	var f = inFile{
		UUID:  strconv.Itoa(s.count),
		Short: short,
		Long:  long,
		Cook:  user,
		Del:   flag,
	}
	s.count += 1
	jsonData, err := json.Marshal(f)
	if err != nil {
		return err
	}

	_, err = s.File.Write(append(jsonData, byte('\n')))
	if err != nil {
		return err
	}

	return nil
}

// GetShortIfHave get short
func (s *FileStorage) GetShortIfHave(user int, path string) (string, error) {
	err := s.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file storage %w", err)
	}
	defer s.File.Close()
	scanner := bufio.NewScanner(s.File)
	for scanner.Scan() {
		line := scanner.Bytes()
		var f inFile
		err := json.Unmarshal(line, &f)
		if err != nil {
			_ = fmt.Errorf("error opening file storage %w", err)
			continue
		}
		if f.Long == path && f.Cook == user {
			return f.Short, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error opening file storage %w", err)
	}
	return "", nil
}

type inFile struct {
	UUID  string `json:"uuid"`
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
	Cook  int    `json:"cook"`
	Del   bool   `json:"is_deleted"`
}

// GetLong get long
func (s *FileStorage) GetLong(user int, short string) (string, error) {
	err := s.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file storage %w", err)
	}
	defer s.File.Close()

	scanner := bufio.NewScanner(s.File)
	for scanner.Scan() {
		line := scanner.Bytes()
		var f inFile
		err := json.Unmarshal(line, &f)
		if err != nil {
			_ = fmt.Errorf("error opening file storage %w", err)
			continue
		}
		if f.Short == short {
			return f.Long, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error opening file storage %w", err)
	}
	return "", nil
}

// GetAll get all
func (s *FileStorage) GetAll(user int) ([]URL, error) {
	err := s.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening file storage %w", err)
	}
	defer s.File.Close()
	var list []URL
	scanner := bufio.NewScanner(s.File)
	for scanner.Scan() {
		line := scanner.Bytes()
		var f inFile
		err := json.Unmarshal(line, &f)
		if err != nil {
			_ = fmt.Errorf("error opening file storage %w", err)
			continue
		}
		url := URL{}
		if f.Cook == user {
			url.Original = f.Long
			url.Short = f.Short
			list = append(list, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error opening file storage %w", err)
	}
	return list, nil
}
