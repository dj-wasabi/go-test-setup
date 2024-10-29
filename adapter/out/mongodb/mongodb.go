package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/port/out"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongodbConnection struct {
	Config  *config.Config
	Logging *slog.Logger
	Client  *mongo.Client
	Context context.Context
	// Cancel  context.CancelFunc
}

func NewMongoDBConnection(c *config.Config) out.OrganisationsDBPort {
	return connectServer(c)
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

func connectServer(c *config.Config) *mongodbConnection {
	logger := logging.Initialize()

	// Not sure why yet, but when commented code is running it will loose connection
	// to MongoDB after a while...
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx := context.Background()
	myConnectString := connectDBString(c)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(myConnectString))
	if err != nil {
		logger.Error(fmt.Sprintf("Mongo DB Connect issue %s", err.Error()))
	}

	con := &mongodbConnection{
		Config:  c,
		Logging: logger,
		Client:  client,
		Context: ctx,
		// Cancel:  cancel,
	}

	return con
}

func (mc *mongodbConnection) pingServer(client *mongo.Client, ctx context.Context) (bool, error) {
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false, err
	}
	return true, nil
}

// func (mc *mongodbConnection) SetupCollection(col string) (*mongo.Collection, *mongo.Client, context.Context, context.CancelFunc) {
// 	_, err := mc.pingServer(mc.Client, mc.Context)
// 	if err != nil {
// 		mc.Logging.Error(fmt.Sprintf("Mongo DB ping issue %s", err.Error()))
// 	}

// 	collection := mc.Client.Database(mc.Config.Database.Dbname).Collection(col)
// 	return collection, client, ctx, cancel
// }

func (mc *mongodbConnection) SetupCollection(col string) *mongo.Collection {
	_, err := mc.pingServer(mc.Client, mc.Context)
	if err != nil {
		// mc.Logging.Error(fmt.Sprintf("Mongo DB ping issue %s", err.Error()))
		con := connectServer(mc.Config)
		collection := con.Client.Database(mc.Config.Database.Dbname).Collection(col)
		return collection
	}

	collection := mc.Client.Database(mc.Config.Database.Dbname).Collection(col)
	return collection
}

func (mc *mongodbConnection) VerifyServer() bool {
	_, _ = mc.pingServer(mc.Client, mc.Context)

	// if err != nil {
	// 	mc.Logging.Error(fmt.Sprintf("Mongo DB ping issue %s", err.Error()))
	// 	return false
	// }

	// defer mc.CloseConnection(mc.Cancel)
	return true
}

// func (mc *mongodbConnection) CloseConnection(cancel context.CancelFunc) {
// 	defer func() {
// 		cancel()
// 		if err := mc.Client.Disconnect(mc.Context); err != nil {
// 			mc.Logging.Error(err.Error())
// 		}
// 	}()
// }
