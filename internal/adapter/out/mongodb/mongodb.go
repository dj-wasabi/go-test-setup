package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

type MongodbRepository struct {
	DB         *mongo.Database
	Collection *mongo.Collection
}

func NewMongodbConnection(c *config.Config) *mongo.Database {
	return connectServer(c)
}

func connectDBString(c *config.Config) string {
	var connectString = "mongodb://"

	if c.Database.Username != "" {
		connectString = connectString + c.Database.Username
	}
	if c.Database.Password != "" {
		connectString = connectString + ":" + c.Database.Password + "@"
	}
	connectString = connectString + c.Database.Hostname + ":" + fmt.Sprintf("%v", c.Database.Port)
	return connectString
}

func connectServer(c *config.Config) *mongo.Database {
	logger := logging.Initialize()

	ctx := context.Background()
	myConnectString := connectDBString(c)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(myConnectString))
	if err != nil {
		logger.Error(fmt.Sprintf("Mongo DB Connect issue %s", err.Error()))
	}
	// Connect to the correct database
	return client.Database(c.Database.Dbname)
}
