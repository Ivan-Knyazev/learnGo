package database

import (
	"go-storage/internal/pkg/models"
	"go-storage/internal/pkg/utils"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(utils.GetEnvs()["DSN"]), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create tables, migrate db
	err = db.AutoMigrate(&models.Data{}, &models.Scalar{}, &models.Slice{}, &models.Dict{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	return db, err
}
