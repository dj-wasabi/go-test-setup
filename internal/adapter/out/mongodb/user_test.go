package mongodb

import (
	"context"
	"os"
	"testing"
	"time"

	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/logging"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_authenticate_GetByName(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	os.Setenv("CONFIGURATION_FILE", "../../../../config.yaml")

	mt.Run("getbyname_result_yes", func(mt *mtest.T) {

		ctx := context.TODO()
		repoUser := NewUserMongoRepo(mt.DB, "users")
		serviceUser := NewUserMongoService(repoUser, logging.Initialize())

		myUser := &out.UserPort{
			ID:        primitive.NewObjectID(),
			Username:  "myusername",
			Password:  "mysecretpassword",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: myUser.ID},
				{Key: "username", Value: myUser.Username},
				{Key: "password", Value: myUser.Password},
				{Key: "enabled", Value: myUser.Enabled},
				{Key: "created_at", Value: myUser.CreatedAt},
				{Key: "updated_at", Value: myUser.UpdatedAt},
			}))

		userData, err := serviceUser.GetByName("myusername", ctx)

		assert.Nil(t, err)
		assert.Equal(t, userData.GetUsername(), "myusername")
		assert.Equal(t, userData.Enabled, true)
	})

	mt.Run("getbyname_result_no", func(mt *mtest.T) {

		ctx := context.TODO()
		repoUser := NewUserMongoRepo(mt.DB, "users")
		serviceUser := NewUserMongoService(repoUser, logging.Initialize())

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch, bson.D{}))
		userData, err := serviceUser.GetByName("myusername", ctx)

		assert.Nil(t, err)
		assert.Equal(t, userData.GetUsername(), "")
	})

}
