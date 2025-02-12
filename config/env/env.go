package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv(key string) string {
	envFiles := []string{
		"../../.env",
		".env",
		filepath.Join(os.Getenv("HOME"), ".env"),
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			break
		}
	}

	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not found", key)
	}
	return value
}
