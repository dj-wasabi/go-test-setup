package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
	"werner-dijkerman.nl/test-setup/pkg/validator"
)

var (
	once           sync.Once
	configInstance *Config
	configError    error
)

var customConfigErrorMessages = map[string]string{
	"Hostname.required": "The field 'DATABASE_HOSTNAME' is required.",
	"Hostname.hostname": "Need to have a proper hostname value.",
	"Port.required":     "The field 'DATABASE_PORT' is required.",
	"Port.numeric":      "Need to have a proper (numeric) port value.",
	"Dbname.required":   "The field 'DATABASE_DBNAME' is required.",
	"Idle.required":     "The field 'HTTP_TIMEOUT_IDLE' is required.",
	"Idle.numeric":      "Need to have a proper (numeric) port value.",
	"Read.required":     "The field 'HTTP_TIMEOUT_READ' is required.",
	"Read.numeric":      "Need to have a proper (numeric) port value.",
	"Write.required":    "The field 'HTTP_TIMEOUT_WRITE' is required.",
	"Write.numeric":     "Need to have a proper (numeric) port value.",
	"Level.oneof":       "Only one of the 'debug', 'info', 'warn', 'warning' or 'error' are allowed.",
}

type Config struct {
	Http     http     `yaml:"http"`
	Database database `yaml:"database"`
	Logging  logging  `yaml:"logging"`
}

type database struct {
	Hostname string `yaml:"hostname" env:"DATABASE_HOSTNAME" validate:"required,hostname"`
	Port     int    `yaml:"port,omitempty" envDefault:"27017" env:"DATABASE_PORT" validate:"required,numeric"`
	Username string `yaml:"username,omitempty" env:"DATABASE_USERNAME"`
	Password string `yaml:"password,omitempty" env:"DATABASE_PASSWORD"`
	Dbname   string `yaml:"dbname" env:"DATABASE_DBNAME" validate:"required"`
}

type http struct {
	Listen  string  `yaml:"listen" env:"HTTP_LISTEN"`
	Logfile string  `yaml:"logfile" env:"HTTP_LOGFILE"`
	Timeout timeout `yaml:"timeout"`
	Cors    cors    `yaml:"cors"`
}

type cors struct {
	Host string `yaml:"host" env:"HTTP_CORS_HOST"`
}

type logging struct {
	Level string `yaml:"level" env:"LOGGING_LEVEL" validate:"oneof=debug info warn warning error"`
}

type timeout struct {
	Idle  int `yaml:"idle" env:"HTTP_TIMEOUT_IDLE" validate:"required,numeric"`
	Read  int `yaml:"read" env:"HTTP_TIMEOUT_READ" validate:"required,numeric"`
	Write int `yaml:"write" env:"HTTP_TIMEOUT_WRITE" validate:"required,numeric"`
}

func loadConfig() (*Config, error) {
	configurationFilePath := os.Getenv("CONFIGURATION_FILE")
	if configurationFilePath == "" {
		configurationFilePath = "config.yaml"
	}
	configurationFilePath = filepath.Clean(configurationFilePath)
	yamlFile, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &configInstance)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&configInstance.Http)
	if err != nil {
		slog.Info(fmt.Sprintf("Error: %v", err))
	}
	err = env.Parse(&configInstance.Http.Timeout)
	if err != nil {
		slog.Info(fmt.Sprintf("Error: %v", err))
	}
	err = env.Parse(&configInstance.Http.Cors)
	if err != nil {
		slog.Info(fmt.Sprintf("Error: %v", err))
	}
	err = env.Parse(&configInstance.Database)
	if err != nil {
		slog.Info(fmt.Sprintf("Error: %v", err))
	}
	err = env.Parse(&configInstance.Logging)
	if err != nil {
		slog.Info(fmt.Sprintf("Error: %v", err))
	}

	err = validator.CheckConfig(*configInstance, customConfigErrorMessages)
	return configInstance, err
}

func ReadConfig() *Config {
	once.Do(func() {
		configInstance, configError = loadConfig()
	})
	if configError != nil {
		panic(fmt.Sprintf("Errors found in configuration: %v", configError))
	}
	return configInstance
}
