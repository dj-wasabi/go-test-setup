package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/adapter/out/mongodb"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["time"] = params.TimeStamp.Format("2006-01-02 15:04:05")
			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["remote_addr"] = params.ClientIP
			log["user_agent"] = params.Request.UserAgent()
			log["referer"] = params.Request.Referer()
			// time endresult is logged in milliseconds: 1000943 == 1 second.
			log["duration"] = time.Duration(params.Latency) / 1000

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func ValidateSecurityScheme(mc *mongodb.MongodbConnection, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	mc.Logging.Info("Getting a validation Security Scheme request")
	clientToken, err := utils.GetBearerToken(mc.Logging, input.RequestValidationInput.Request)
	if err != nil {
		mc.Logging.Error(fmt.Sprintf("%v", err.Error()))
		return err
	}

	claims, err := utils.ValidateToken(mc.Logging, clientToken)
	if err != nil {
		myError := model.GetError("AUTH001")
		return errors.New(myError.Message)
	}

	user, _ := mc.GetByName(claims.Username, ctx)
	if clientToken == user.Token {
		return nil
	} else {
		myError := model.GetError("AUTH002")
		return errors.New(myError.Message)
	}
}

// do some magic with mongodb database, getting token from user
