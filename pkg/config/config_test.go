package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadconfig_configfile(t *testing.T) {
	config, err := loadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, config.Logging.Level, "info")
	assert.Equal(t, config.Database.Port, 27017)
}

func Test_loadconfig_error(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "DEBUG")
	_, err := loadConfig()

	assert.Equal(t, err.Error(), "Only one of the 'debug', 'info', 'warn', 'warning' or 'error' are allowed.")
}

func Test_loadconfig_override(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "debug")
	os.Setenv("DATABASE_PORT", "27019")
	os.Setenv("HTTP_TIMEOUT_IDLE", "1")

	config, err := loadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, config.Logging.Level, "debug")
	assert.Equal(t, config.Database.Port, 27019)
	assert.Equal(t, config.Http.Timeout.Idle, 1)
}

func Test_loadconfig_ignore_error(t *testing.T) {
	os.Setenv("HTTP_TIMEOUT_IDLE", "pizza")

	config, err := loadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, config.Http.Timeout.Idle, 10)
}
