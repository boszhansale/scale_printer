package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"test/internal/logger"
)

type Config struct {
	WeightAddress string `env:"WEIGHT_ADDRESS"`
	PrinterName   string `env:"PRINTER_NAME"`
}

func NewConfig() *Config {
	path := "./.env"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		logger.Error("failed to read config: %v", err)
	}
	return &cfg
}
