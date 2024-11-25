package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

// adapter

type userService struct {
	logging    *slog.Logger
	repository *MongodbRepository
}

func NewUserMongoRepo(con *mongo.Database, collection string) *MongodbRepository {
	return &MongodbRepository{
		DB:         con,
		Collection: con.Collection(collection),
	}
}

func NewUserMongoService(repo *MongodbRepository, log *slog.Logger) out.PortUser {
	return &userService{
		logging:    log,
		repository: repo,
	}
}

// asasasasasa (These come from port/out/(interface))
func (uc *userService) Create(ctx context.Context, user *out.UserPort) (string, *model.Error) {
	newUser, _ := model.NewUser(user.GetUsername(), user.GetPassword(), user.GetRole(), user.GetEnabled())

	uc.logging.Debug(fmt.Sprintf("Creating account with username '%v'", user.GetUsername()))
	add, err := uc.repository.Collection.InsertOne(ctx, newUser)

	if err != nil {
		var write_exc mongo.WriteException
		if !errors.As(err, &write_exc) {
			uc.logging.Error(fmt.Sprintf("%v", err))
			return "", model.GetError("UNKNOWN")
		}

		if write_exc.HasErrorCodeWithMessage(11000, "index: unique_username_idx") {
			uc.logging.Error(fmt.Sprintf("User '%v' already exist, unique index violation.", user.GetUsername()))
			return "", model.GetError("USR0001")
		}
	}

	oid, _ := add.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func (uc *userService) GetByName(username string, ctx context.Context) (*out.UserPort, error) {
	uc.logging.Debug(fmt.Sprintf("About to Create Organisations %v", username))

	result := uc.repository.Collection.FindOne(ctx, bson.M{"username": username})
	if result.Err() == mongo.ErrNoDocuments {
		uc.logging.Info(fmt.Sprintf("User '%v' not found.", username))
		return nil, result.Err()
	}

	user := new(*out.UserPort)
	err := result.Decode(&user)
	if err != nil {
		uc.logging.Error("error")
	}

	return *user, nil
}
