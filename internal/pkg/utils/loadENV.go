package utils

import (
	"go-storage/internal/pkg/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func GetEnvs(s *storage.Storage) map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	evns := make(map[string]string)
	evns["path"] = os.Getenv("JSON_PATH")
	evns["port"] = os.Getenv("PORT")

	s.Logger.Info("ENV was loaded", zap.String("JSON_PATH", evns["path"]), zap.Any("PORT", evns["port"]))
	defer s.Logger.Sync()

	return evns
}
