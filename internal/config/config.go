package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port   string `envconfig:"PORT" required:"true"`
	DBHost string `envconfig:"DB_HOST" required:"true"`
	DBPort string `envconfig:"DB_PORT" required:"true"`
	DBUser string `envconfig:"DB_USER" required:"true"`
	DBPass string `envconfig:"DB_PASSWORD" required:"true"`
	DBName string `envconfig:"DB_NAME" required:"true"`
}

func Load() *Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	return &cfg
}
