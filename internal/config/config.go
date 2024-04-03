package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	CategoriesApi    string `env:"CATEGORIES_API" envDefault:""`
	ProductsApi      string `env:"PRODUCTS_API" envDefault:""`
	ProductCreateApi string `env:"PRODUCT_CREATE_API" envDefault:""`
	WeightAddress    string `env:"WEIGHT_ADDRESS"`
	PrinterName      string `env:"PRINTER_ADDRESS"`
}

func NewConfig() *Config {
	path := "./.env"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
