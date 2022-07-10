package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err)
	}
	log.Println("Loaded .env file")
}

func Get(key string) string {
	return os.Getenv(key)
}
