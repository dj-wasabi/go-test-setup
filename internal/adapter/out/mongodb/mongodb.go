package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongodbConnection struct {
	Config  *config.Config
	Logging *slog.Logger
	Client  *mongo.Client
	Context context.Context
}

func NewMongodbConnection(c *config.Config) (*MongodbConnection, *MongodbConnection) {
	return connectServer(c), connectServer(c)
}

func connectDBString(mc *config.Config) string {
	var connectString = "mongodb://"

	if mc.Database.Username != "" {
		connectString = connectString + mc.Database.Username
	}
	if mc.Database.Password != "" {
		connectString = connectString + ":" + mc.Database.Password + "@"
	}
	connectString = connectString + mc.Database.Hostname + ":" + fmt.Sprintf("%v", mc.Database.Port)
	return connectString
}

func connectServer(c *config.Config) *MongodbConnection {
	logger := logging.Initialize()

	ctx := context.Background()
	myConnectString := connectDBString(c)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(myConnectString))
	if err != nil {
		logger.Error(fmt.Sprintf("Mongo DB Connect issue %s", err.Error()))
	}

	con := &MongodbConnection{
		Config:  c,
		Logging: logger,
		Client:  client,
		Context: ctx,
	}

	return con
}

func (mc *MongodbConnection) pingServer(client *mongo.Client, ctx context.Context) (bool, error) {
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (mc *MongodbConnection) SetupCollection(col string) *mongo.Collection {
	_, err := mc.pingServer(mc.Client, mc.Context)
	if err != nil {
		mc.Logging.Error(fmt.Sprintf("Mongo DB ping issue %s", err.Error()))
		con := connectServer(mc.Config)
		mc.Logging.Error("Trying to get another connection")
		collection := con.Client.Database(mc.Config.Database.Dbname).Collection(col)
		return collection
	}

	collection := mc.Client.Database(mc.Config.Database.Dbname).Collection(col)
	return collection
}

func (mc *MongodbConnection) VerifyServer() bool {
	_, _ = mc.pingServer(mc.Client, mc.Context)
	return true
}
