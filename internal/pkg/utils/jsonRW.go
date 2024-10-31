package utils

import (
	"encoding/json"
	"go-storage/internal/pkg/storage"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func getFilePath() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	path := os.Getenv("JSON_PATH")
	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	filePath := filepath.Join(dir, filename)
	return filePath
}

func writeAtomic(data []byte) error {
	filepath := getFilePath()
	tmpFilepath := getFilePath() + ".tmp"
	err := os.WriteFile(tmpFilepath, data, 0744)
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(tmpFilepath)
	}()

	return os.Rename(tmpFilepath, filepath)
}

func ReadFromFile(s *storage.Storage) error {
	filePath := getFilePath()
	fromFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var data storage.JsonStorage
	err = json.Unmarshal(fromFile, &data)
	if err != nil {
		return err
	}

	s.LoadData(data)
	log.Println(*s)
	return nil
}

func WriteToFile(s *storage.Storage) error {
	jsonStorage := s.ExportData()
	log.Println(jsonStorage)
	data, err := json.MarshalIndent(jsonStorage, "", "\t")
	if err != nil {
		return err
	}

	err = writeAtomic(data)
	if err != nil {
		return err
	}

	return nil
}
