package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type FileStorage struct {
	Path  string
	File  *os.File
	count int
}

func (s *FileStorage) CheckConnection() error {
	if s.File == nil {
		return fmt.Errorf("file not open")
	}
	return nil

}

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

func (s *FileStorage) Open() error {
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	s.File = file
	return nil
}

func (s *FileStorage) Close() error {
	return s.File.Close()
}

func (s *FileStorage) Save(short string, long string) error {

	var f = inFile{
		UUID:  strconv.Itoa(s.count),
		Short: short,
		Long:  long,
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

func (s *FileStorage) GetShortIfHave(path string) (string, error) {
	err := s.Open()
	if err != nil {
		fmt.Println("Not open")
		return "", err
	}
	defer s.File.Close()

	scanner := bufio.NewScanner(s.File)
	for scanner.Scan() {
		line := scanner.Bytes()
		var f inFile
		err := json.Unmarshal(line, &f)
		if err != nil {
			fmt.Println("Bad json in File", err)
			continue
		}
		if f.Long == path {
			return f.Short, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Bad file")
		return "", err
	}
	return "", nil
}

type inFile struct {
	UUID  string `json:"uuid"`
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

func (s *FileStorage) GetLong(short string) (string, error) {
	err := s.Open()
	if err != nil {
		fmt.Println("Not open")
		return "", err
	}
	defer s.File.Close()

	scanner := bufio.NewScanner(s.File)
	for scanner.Scan() {
		line := scanner.Bytes()
		var f inFile
		err := json.Unmarshal(line, &f)
		if err != nil {
			fmt.Println("Bad json in File", err)
			continue
		}
		if f.Short == short {
			return f.Long, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Bad file")
		return "", err
	}
	return "", nil
}
