package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectID              string
	RabbitMQURL            string
	CollectionName         string
	ApplicationCredentials string
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{
		ProjectID:              os.Getenv("GOOGLE_CLOUD_PROJECT"),
		RabbitMQURL:            os.Getenv("RABBITMQ_URL"),
		CollectionName:         os.Getenv("FIRESTORE_COLLECTION"),
		ApplicationCredentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}

	// TODO: env validadtion

	return cfg, nil
}
