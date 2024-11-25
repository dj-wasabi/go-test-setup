package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
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
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

var (
	handler             *ApiHandler
	router              *gin.Engine
	mongoTest           *mtest.T
	domainService       in.ApiUseCases
	repoUser            *mongodb.MongodbRepository
	serviceUser         out.PortUser
	serviceOrganisation out.PortOrganisation
	authToken           *model.AuthenticatePostResponse
	authRequest         model.AuthenticatePostRequest
	authError           model.Error
	myUser              out.UserPort
	token               string
)

func prepareTest(t *testing.T) {
	mongoTest = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	os.Setenv("LOGFILE_PATH", "../../../../../config.yaml")
}

func prepareMongoDB(mt *mtest.T, l *slog.Logger) {
	repoUser = mongodb.NewUserMongoRepo(mt.DB, "users")
	serviceUser = mongodb.NewUserMongoService(repoUser, l)
	_ = config.ReadConfig()
	token, _ = utils.GenerateToken("myusername", "admin")

	myUser = out.UserPort{
		ID:        primitive.NewObjectID(),
		Username:  "myusername",
		Password:  "$2a$14$flIjKE7ywigEp8c5.7TFru8OKiTXMz0TG21TmwL8jfmnvMOHvj0Oi",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      "admin",
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
}

func prepareGin() {
	router = gin.New()
	gin.SetMode(gin.TestMode)

	domainService = services.NewdomainServices(serviceOrganisation, serviceUser)
	handler = NewApiService(domainService)
	RegisterHandlers(router, handler)
}

func Test_Authenticatelogin_Ok(t *testing.T) {
	prepareTest(t)

	mongoTest.Run("authenticate", func(mt *mtest.T) {

		prepareMongoDB(mt, logging.Initialize())
		prepareGin()

		authRequest = model.AuthenticatePostRequest{
			Username: myUser.Username,
			Password: "mysecretpassword",
		}
		b, _ := json.Marshal(authRequest)
		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewReader(b))
		req.Header.Set("Accept", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		err := json.NewDecoder(rec.Body).Decode(&authToken)

		assert.NoError(t, err)
		assert.NotEqual(t, authToken.Token, "")
	})

}

func Test_Authenticatelogin_NotOk(t *testing.T) {
	prepareTest(t)

	mongoTest.Run("authenticate", func(mt *mtest.T) {

		prepareMongoDB(mt, logging.Initialize())
		prepareGin()

		authRequest = model.AuthenticatePostRequest{
			Username: myUser.Username,
			Password: "mysecretpasswor",
		}
		b, _ := json.Marshal(authRequest)
		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewReader(b))
		req.Header.Set("Accept", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		err := json.NewDecoder(rec.Body).Decode(&authError)

		assert.NoError(t, err)
		assert.Equal(t, authError.Message, "Invalid username/password combination")
	})

}
