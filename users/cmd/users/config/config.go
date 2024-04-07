package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env     string        `yaml:"env"`
	Storage StorageConfig `yaml:"storage"`
	Server  ServerConfig  `yaml:"server"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `env:"DB_PASSWORD"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	var cfg Config

	configPath := flag.String("config", "", "config path")
	flag.Parse()

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		return nil, fmt.Errorf("can't load config: %w", err)
	}

	return &cfg, nil
}
