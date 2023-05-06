package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramApiToken string `env:"TELEGRAM_API_TOKEN"`
}

func New() (*Config, error) {
	err := loadEnv()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadEnv() error {
	err := godotenv.Load(os.Getenv("ENV_FILE"))
	if err != nil {
		return err
	}
	return nil
}
