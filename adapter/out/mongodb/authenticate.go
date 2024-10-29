package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/port/out"
)

// asasasasasa (These come from port/out/(interface))
func (mc *mongodbConnection) GetUserByName(username string, ctx context.Context) (*out.User, error) {
	var mdbCollection string = "user"
	mc.Logging.Debug(fmt.Sprintf("About to Create Organisations %v", username))

	coll := mc.SetupCollection(mdbCollection)

	result := coll.FindOne(mc.Context, bson.M{"username": username})
	mc.Logging.Info(fmt.Sprintf("%v", result.Err()))
	if result.Err() == mongo.ErrNoDocuments {
		mc.Logging.Info(fmt.Sprintf("User '%v' not found.", username))
		return nil, result.Err()
	}

	user := new(*out.User)
	err := result.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	return *user, nil
}
