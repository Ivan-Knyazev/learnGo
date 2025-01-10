package storage

import (
	"log"

	"gorm.io/gorm"
)

type Saver interface {
	SaveData(storage Storage) error
	ReadData(storage Storage) error
}

type JSONSaver struct {
	filePath string
}

func CreateJSONSaver(filePath string) *JSONSaver {
	return &JSONSaver{filePath: filePath}
}

func (s *JSONSaver) ReadData(storage Storage) error {
	err := ReadFromFile(storage, s.filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *JSONSaver) SaveData(storage Storage) error {
	err := WriteToFile(storage, s.filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type DBSaver struct {
	db *gorm.DB
}

func CreateDBSaver(db *gorm.DB) *DBSaver {
	return &DBSaver{db: db}
}

func (s *DBSaver) ReadData(storage Storage) error {
	err := ReadFromDB(storage, s.db)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *DBSaver) SaveData(storage Storage) error {
	err := WriteToDB(storage, s.db)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
