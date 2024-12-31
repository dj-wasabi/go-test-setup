package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"werner-dijkerman.nl/test-setup/pkg/config"
)

func Test_connectDBString(t *testing.T) {
	config := &config.Config{}
	config.Database.Username = "myuser"
	config.Database.Password = "password"
	config.Database.Hostname = "localhost"
	config.Database.Port = 12311

	connectString := connectDBString(config)
	assert.Equal(t, connectString, "mongodb://myuser:password@localhost:12311")
}

func Test_connectDBString_no_creds(t *testing.T) {
	config := &config.Config{}
	config.Database.Hostname = "localhost"
	config.Database.Port = 12311

	connectString := connectDBString(config)
	assert.Equal(t, connectString, "mongodb://localhost:12311")
}
