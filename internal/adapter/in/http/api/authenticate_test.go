package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"werner-dijkerman.nl/test-setup/internal/adapter/out/mongodb"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/domain/services"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func Test_Authenticatelogin_Ok(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	os.Setenv("LOGFILE_PATH", "../../../../../config.yaml")

	mt.Run("authenticate", func(mt *mtest.T) {

		repoUser := mongodb.NewUserMongoRepo(mt.DB, "users")
		serviceUser := mongodb.NewUserMongoService(repoUser, logging.Initialize())
		token, _ := utils.GenerateToken("myusername")

		myUser := &out.UserPort{
			ID:        primitive.NewObjectID(),
			Username:  "myusername",
			Password:  "$2a$14$flIjKE7ywigEp8c5.7TFru8OKiTXMz0TG21TmwL8jfmnvMOHvj0Oi",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Token:     token,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: myUser.ID},
				{Key: "username", Value: myUser.Username},
				{Key: "password", Value: myUser.Password},
				{Key: "enabled", Value: myUser.Enabled},
				{Key: "created_at", Value: myUser.CreatedAt},
				{Key: "updated_at", Value: myUser.UpdatedAt},
				{Key: "token", Value: myUser.Token},
			}))
		_ = config.ReadConfig()

		r := gin.New()
		gin.SetMode(gin.TestMode)
		var serviceOrganisation out.PortOrganisation

		ds := services.NewdomainServices(serviceOrganisation, serviceUser)
		h := NewApiService(ds)
		RegisterHandlers(r, h)

		body := model.AuthenticationRequest{
			Username: myUser.Username,
			Password: "mysecretpassword",
		}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(b))
		req.Header.Set("Accept", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var authToken *model.AuthenticationToken
		err := json.NewDecoder(rec.Body).Decode(&authToken)

		assert.NoError(t, err)
		assert.NotEqual(t, authToken.Token, "")
	})

}

func Test_Authenticatelogin_NotOk(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	os.Setenv("LOGFILE_PATH", "../../../../../config.yaml")

	mt.Run("authenticate", func(mt *mtest.T) {

		repoUser := mongodb.NewUserMongoRepo(mt.DB, "users")
		serviceUser := mongodb.NewUserMongoService(repoUser, logging.Initialize())
		token, _ := utils.GenerateToken("myusername")

		myUser := &out.UserPort{
			ID:        primitive.NewObjectID(),
			Username:  "myusername",
			Password:  "$2a$14$flIjKE7ywigEp8c5.7TFru8OKiTXMz0TG21TmwL8jfmnvMOHvj0Oi",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Token:     token,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: myUser.ID},
				{Key: "username", Value: myUser.Username},
				{Key: "password", Value: myUser.Password},
				{Key: "enabled", Value: myUser.Enabled},
				{Key: "created_at", Value: myUser.CreatedAt},
				{Key: "updated_at", Value: myUser.UpdatedAt},
				{Key: "token", Value: myUser.Token},
			}))
		_ = config.ReadConfig()

		r := gin.New()
		gin.SetMode(gin.TestMode)
		var serviceOrganisation out.PortOrganisation

		ds := services.NewdomainServices(serviceOrganisation, serviceUser)
		h := NewApiService(ds)
		RegisterHandlers(r, h)

		body := model.AuthenticationRequest{
			Username: myUser.Username,
			Password: "mysecretpasswor",
		}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(b))
		req.Header.Set("Accept", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var authToken *model.Error
		err := json.NewDecoder(rec.Body).Decode(&authToken)

		assert.NoError(t, err)
		assert.Equal(t, authToken.Message, "Invalid username/password combination")
	})

}
