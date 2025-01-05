package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvs() map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	envs := make(map[string]string)
	envs["path"] = os.Getenv("JSON_PATH")
	envs["port"] = os.Getenv("PORT")
	envs["interval"] = os.Getenv("INTERVAL")
	envs["timeout"] = os.Getenv("SHUTDOWN_TIMEOUT")

	// Create DSN (Data Source Name)
	POSTGRES_DB := os.Getenv("POSTGRES_DB")
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	envs["DSN"] = fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT)

	return envs
}
