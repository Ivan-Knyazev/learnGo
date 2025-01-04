package utils

import (
	"fmt"
	"go-storage/internal/pkg/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvs(s storage.Storage) map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	envs := make(map[string]string)
	envs["path"] = os.Getenv("JSON_PATH")
	envs["port"] = os.Getenv("PORT")
	envs["interval"] = os.Getenv("INTERVAL")
	envs["timeout"] = os.Getenv("SHUTDOWN_TIMEOUT")

	logString := fmt.Sprintf("ENV was loaded, JSON_PATH=%s, PORT=%s", envs["path"], envs["port"])
	defer s.WriteLog(logString)

	return envs
}
