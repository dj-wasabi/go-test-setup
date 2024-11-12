package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"werner-dijkerman.nl/test-setup/pkg/validator"
)

func Test_loadconfig_configfile(t *testing.T) {
	config, err := loadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, config.Logging.Level, "INFO")
	assert.Equal(t, config.Database.Port, 27017)

}

func Test_loadconfig_override(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "DEBUG")
	os.Setenv("DATABASE_PORT", "27019")
	os.Setenv("HTTP_TIMEOUT_IDLE", "1")

	config, err := loadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, config.Logging.Level, "DEBUG")
	assert.Equal(t, config.Database.Port, 27019)
	assert.Equal(t, config.Http.Timeout.Idle, 1)
}

func Test_validateConfig_be_true(t *testing.T) {
	os.Setenv("HTTP_TIMEOUT_IDLE", "1")

	config, _ := loadConfig()
	v := validator.New()
	okConfig := validateConfig(v, config)
	assert.True(t, okConfig)
}

func Test_validateConfig_incorrect_value(t *testing.T) {
	os.Setenv("HTTP_TIMEOUT_IDLE", "0")
	os.Setenv("HTTP_TIMEOUT_WRITE", "0")

	config, _ := loadConfig()
	v := validator.New()
	okConfig := validateConfig(v, config)
	assert.False(t, okConfig)
}
