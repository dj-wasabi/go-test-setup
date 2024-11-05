package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

// asasasasasa (These come from port/out/(interface))
// func (mc *mongodbConnection) Create(ctx context.Context, user *out.UserPort) string {
// 	var mdbCollection string = "users"

// 	newUser := model.NewUser(user.GetUsername(), user.GetPassword(), user.GetEnabled(), user.GetRoles())
// 	coll := mc.SetupCollection(mdbCollection)

// 	result := coll.FindOne(ctx, bson.M{"username": user.GetUsername()})
// 	mc.Logging.Debug(fmt.Sprintf("Creating account with username '%v'", user.GetUsername()))
// 	if result.Err() == mongo.ErrNoDocuments {
// 		add, _ := coll.InsertOne(ctx, newUser)
// 		oid, ok := add.InsertedID.(primitive.ObjectID)
// 		if !ok {
// 			mc.Logging.Error(fmt.Sprintf("Error while getting the object id for user %v.", user.GetUsername()))
// 		}
// 		mc.Logging.Debug(fmt.Sprintf("Account '%v' created with id '%v'", user.GetUsername(), oid.Hex()))
// 		return oid.Hex()
// 	} else {
// 		return "duplicate"
// 	}
// }

func (mc *mongodbConnection) Create(ctx context.Context, user *out.UserPort) (string, *model.Error) {
	var mdbCollection string = "users"

	newUser := model.NewUser(user.GetUsername(), user.GetPassword(), user.GetEnabled(), user.GetRoles())
	coll := mc.SetupCollection(mdbCollection)

	mc.Logging.Debug(fmt.Sprintf("Creating account with username '%v'", user.GetUsername()))
	add, err := coll.InsertOne(ctx, newUser)

	if err != nil {
		var write_exc mongo.WriteException
		if !errors.As(err, &write_exc) {
			mc.Logging.Error(fmt.Sprintf("%v", err))
			return "", model.GetError("UNKNOWN")
		}

		if write_exc.HasErrorCodeWithMessage(11000, "index: unique_username_idx") {
			mc.Logging.Error(fmt.Sprintf("User '%v' already exist, unique index violation.", user.GetUsername()))
			return "", model.GetError("USR0001")
		}
	}

	oid, _ := add.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}
