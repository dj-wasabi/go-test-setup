package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"werner-dijkerman.nl/test-setup/pkg/utils"
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

	uc.logging.Debug("log_id", utils.GetLogId(ctx), fmt.Sprintf("Updating token for username %v", username))
	timeStart := time.Now()
	_, err := uc.repository.Collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	timeEnd := float64(time.Since(timeStart).Seconds())

	if err != nil {
		mongodb_user_tokens.WithLabelValues("failure").Observe(timeEnd)
		return false
	}
	mongodb_user_tokens.WithLabelValues("successful").Observe(timeEnd)
	return true
}
