package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type FileStorage struct {
	Path string
	File *os.File
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{Path: path}
}

func (s *FileStorage) Open() error {
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	s.File = file
	return nil
}

func (s *FileStorage) SaveInFile(short string, long string) error {
	err := s.Open()
	if err != nil {
		return err
	}
	type inFile struct {
		Uuid  string `json:"uuid"`
		Short string `json:"short_url"`
		Long  string `json:"original_url"`
	}
	var f inFile
	count := 1
	f.Uuid = strconv.Itoa(count)
	count += 1
	f.Short = short
	f.Long = long
	writer := bufio.NewWriter(s.File)
	defer s.File.Close()
	jsonData, err := json.Marshal(f)

	_, err = writer.WriteString(string(jsonData) + "\n")
	if err != nil {
		fmt.Println("Not write")
		return err
	}
	writer.Flush()
	return nil
}

func (s *FileStorage) GetLong(short string) (string, error) {
	err := s.Open()
	if err != nil {
		fmt.Println("Not open")
		return "", err
	}
	defer s.File.Close()
	type inFile struct {
		Uuid  string `json:"uuid"`
		Short string `json:"short_url"`
		Long  string `json:"original_url"`
	}
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
			return short, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Bad file")
		return "", err
	}
	return "", nil
}
