package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	Env          string            `yaml:"env" env-default:"dev"`
	Application  ApplicationConfig `yaml:"application"`
	PublicServer ServerConfig      `yaml:"public_server"`
	AdminServer  ServerConfig      `yaml:"admin_server"`
	Storage      StorageConfig     `yaml:"storage"`
}

type ApplicationConfig struct {
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout"`
	App                     string        `yaml:"application"`
	SaltValue               string        `yaml:"salt_value"`
}

type ServerConfig struct {
	Enable       bool   `yaml:"enable"`
	Endpoint     string `yaml:"endpoint"`
	Port         int    `yaml:"port" env:"PORT"`
	JwtTokenSalt string `env:"JWT_TOKEN_SALT"`
}

type StorageConfig struct {
	EnableMock            bool          `yaml:"enable_mock"`
	Hosts                 []string      `yaml:"hosts"`
	Port                  int           `yaml:"port"`
	Database              string        `yaml:"database"`
	Username              string        `yaml:"username"`
	Password              string        `yaml:"password" env:"DB_PASSWORD"`
	SSLMode               string        `yaml:"ssl_mode"`
	ConnectionAttempts    int           `yaml:"connection_attempts"`
	InitializationTimeout time.Duration `yaml:"initialization_timeout"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	cfg := NewDefaultConfig()

	configPath := flag.String("config", "", "config path")
	flag.Parse()

	if err := cleanenv.ReadConfig(*configPath, cfg); err != nil {
		return nil, fmt.Errorf("can't load config: %w", err)
	}

	return cfg, nil
}
