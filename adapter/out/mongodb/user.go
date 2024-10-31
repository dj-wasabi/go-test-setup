package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/domain/model"
)

// asasasasasa (These come from port/out/(interface))
func (mc *mongodbConnection) Create(ctx context.Context, user *model.User) error {
	var mdbCollection string = "users"
	username := user.GetUsername()
	mc.Logging.Info(fmt.Sprintf("About to create a user %v", username))
	mc.Logging.Info(fmt.Sprintf("With password %v", user.GetPassword()))

	newUser := model.NewUser(user.GetUsername(), user.GetPassword(), user.GetEnabled(), user.GetRoles())
	coll := mc.SetupCollection(mdbCollection)

	result := coll.FindOne(mc.Context, bson.M{"username": username})
	mc.Logging.Info(fmt.Sprintf("%v", result.Err()))
	if result.Err() == mongo.ErrNoDocuments {
		_, _ = coll.InsertOne(mc.Context, newUser)
		// mc.Logging.Info(fmt.Sprintf("Created with %v", pizza.InsertedID))
	}

	return nil
}
