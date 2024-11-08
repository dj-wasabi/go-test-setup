package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

// asasasasasa (These come from port/out/(interface))
func (mc *MongodbConnection) GetByName(username string, ctx context.Context) (*out.UserPort, error) {
	var mdbCollection string = "users"
	mc.Logging.Debug(fmt.Sprintf("About to Create Organisations %v", username))

	coll := mc.SetupCollection(mdbCollection)

	result := coll.FindOne(mc.Context, bson.M{"username": username})
	mc.Logging.Info(fmt.Sprintf("%v", result.Err()))
	if result.Err() == mongo.ErrNoDocuments {
		mc.Logging.Info(fmt.Sprintf("User '%v' not found.", username))
		return nil, result.Err()
	}

	user := new(*out.UserPort)
	err := result.Decode(&user)
	if err != nil {
		mc.Logging.Error("error")
	}

	return *user, nil
}

func (mc *MongodbConnection) UpdateToken(ctx context.Context, token, username string) bool {
	var mdbCollection string = "users"
	var updateObj primitive.D
	coll := mc.SetupCollection(mdbCollection)

	updateObj = append(updateObj, bson.E{"token", token})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	upsert := true
	filter := bson.M{"username": username}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := coll.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		return false
	}
	return true
}
