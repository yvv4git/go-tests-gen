package config

import (
	"fmt"

	stdLog "log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel string `toml:"log_level" env:"LOG_LEVEL" default:"info"`
	LLM      LLM    `toml:"llm"`
}

func Load(path string, cfg *Config) error {
	if err := godotenv.Load(); err != nil {
		stdLog.Println("Failed read .env file")
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return nil
}
