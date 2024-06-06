package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	prod := os.Getenv("PROD")

	if prod == "TRUE" {
		return nil
	}

	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}