package utils

import (
	"encoding/json"
	"go-storage/internal/pkg/storage"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func getFilePath(file string) string {
	dir := filepath.Dir(file)
	filename := filepath.Base(file)

	filePath := filepath.Join(dir, filename)
	return filePath
}

func writeAtomic(data []byte, file string) error {
	filepath := getFilePath(file)
	tmpFilepath := filepath + ".tmp"
	err := os.WriteFile(tmpFilepath, data, 0744)
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(tmpFilepath)
	}()

	return os.Rename(tmpFilepath, filepath)
}

func ReadFromFile(s *storage.Storage, file string) error {
	filePath := getFilePath(file)
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

	s.Logger.Info("STATE was readed from JSON", zap.Any("data", data))
	defer s.Logger.Sync()

	return nil
}

func WriteToFile(s *storage.Storage, file string) error {
	jsonStorage := s.ExportData()
	data, err := json.MarshalIndent(jsonStorage, "", "\t")
	if err != nil {
		return err
	}

	err = writeAtomic(data, file)
	if err != nil {
		return err
	}

	s.Logger.Info("STATE was writed to JSON", zap.Any("data", jsonStorage))
	defer s.Logger.Sync()

	return nil
}
