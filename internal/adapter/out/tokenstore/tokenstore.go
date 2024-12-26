package tokenstore

import (
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

type tokenstoreService struct {
	logging *slog.Logger
	client  *redis.Client
}

func NewTokenstoreService(token *redis.Client, log *slog.Logger) out.PortStoreInterface {
	return &tokenstoreService{
		logging: log,
		client:  token,
	}
}

func NewTokenstoreConnection(c *config.Config) *redis.Client {
	return connectServer(c)
}

func connectTokenStoreString(c *config.Config) string {
	var connectString = "redis://"

	if c.Tokentstore.Username != "" {
		connectString = connectString + c.Tokentstore.Username
	}
	if c.Tokentstore.Password != "" {
		connectString = connectString + ":" + c.Tokentstore.Password + "@"
	}

	connectString = connectString + c.Tokentstore.Hostname + ":" + fmt.Sprintf("%v", c.Tokentstore.Port)
	connectString = connectString + "/" + c.Tokentstore.Dbname + "?protocol=" + fmt.Sprintf("%v", c.Tokentstore.Protocol)
	return connectString

}

func connectServer(c *config.Config) *redis.Client {
	logger := logging.Initialize()
	url := connectTokenStoreString(c)
	TokenStore, err := redis.ParseURL(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while connecting to redis: %v", err))
	}

	return redis.NewClient(TokenStore)
}
