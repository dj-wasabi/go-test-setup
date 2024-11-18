package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// asasasasasa (These come from port/out/(interface))

func (uc *userService) UpdateToken(ctx context.Context, token, username string) bool {
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", token})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	upsert := true
	filter := bson.M{"username": username}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := uc.repository.Collection.UpdateOne(
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
