package config

import "time"

const (
	defaultAppName   = "avito-tech"
	defaultSaltValue = "saltValue"
)

func NewDefaultConfig() *Config {
	return &Config{
		Application: ApplicationConfig{
			GracefulShutdownTimeout: 15 * time.Second,
			App:                     defaultAppName,
			SaltValue:               defaultSaltValue,
		},
		PublicServer: ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		AdminServer: ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		Storage: StorageConfig{
			EnableMock:            false,
			Hosts:                 []string{},
			Port:                  0,
			Database:              "",
			Username:              "",
			Password:              "",
			SSLMode:               "",
			ConnectionAttempts:    0,
			InitializationTimeout: 5 * time.Second,
		},
	}
}
