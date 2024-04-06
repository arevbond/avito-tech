package config

type Config struct {
	Env string `yaml:"env"`
}

type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `env:"DB_PASSWORD"`
}
