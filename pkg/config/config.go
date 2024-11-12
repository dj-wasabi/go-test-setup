package config

import (
	"fmt"
	"os"
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

type Config struct {
	Http     http     `yaml:"http"`
	Database database `yaml:"database"`
	Logging  logging  `yaml:"logging"`
}

type database struct {
	Hostname string `yaml:"hostname" env:"DATABASE_HOSTNAME"`
	Port     int    `yaml:"port,omitempty" envDefault:"27017" env:"DATABASE_PORT"`
	Username string `yaml:"username,omitempty" env:"DATABASE_USERNAME"`
	Password string `yaml:"password,omitempty" env:"DATABASE_PASSWORD"`
	Dbname   string `yaml:"dbname" env:"DATABASE_DBNAME"`
}

type http struct {
	Listen  string  `yaml:"listen" env:"HTTP_LISTEN"`
	Timeout timeout `yaml:"timeout"`
	Cors    cors    `yaml:"cors"`
}

type cors struct {
	Host string `yaml:"host" env:"HTTP_CORS_HOST"`
}

type logging struct {
	Level string `yaml:"level,omitempty" envDefault:"INFO" env:"LOGGING_LEVEL"`
}

type timeout struct {
	Idle  int `yaml:"idle" env:"HTTP_TIMEOUT_IDLE"`
	Read  int `yaml:"read" env:"HTTP_TIMEOUT_READ"`
	Write int `yaml:"write" env:"HTTP_TIMEOUT_WRITE"`
}

func loadConfig() (*Config, error) {
	logfilePath := os.Getenv("LOGFILE_PATH")
	if logfilePath == "" {
		logfilePath = "config.yaml"
	}
	yamlFile, err := os.ReadFile(logfilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &configInstance)
	if err != nil {
		return nil, err
	}

	_ = env.Parse(&configInstance.Http)
	_ = env.Parse(&configInstance.Http.Timeout)
	_ = env.Parse(&configInstance.Http.Cors)
	_ = env.Parse(&configInstance.Database)
	_ = env.Parse(&configInstance.Logging)

	return configInstance, nil
}

func ReadConfig() *Config {
	once.Do(func() {
		configInstance, configError = loadConfig()
	})
	checkConfig(configError)
	return configInstance
}

func checkConfig(configError error) bool {

	v := validator.New()
	if configError != nil {
		fmt.Println(configError.Error())
	}
	okConfig := validateConfig(v, configInstance)
	if !okConfig {
		for k, v := range v.Errors {
			fmt.Println(k, "value", v)
		}
		os.Exit(1)
	}
	return true
}

func validateConfig(v *validator.Validator, config *Config) bool {

	v.Check(config.Http.Listen != "", "http.listen", "must be provided")
	v.Check(config.Http.Timeout.Idle >= 1, "http.timeout.idle", fmt.Sprintf("must be greater than 1, has %v", config.Http.Timeout.Idle))
	v.Check(config.Http.Timeout.Read >= 1, "http.timeout.read", "must be greater than 1")
	v.Check(config.Http.Timeout.Write >= 1, "http.timeout.write", "must be greater than 1")

	return v.Valid()
}
