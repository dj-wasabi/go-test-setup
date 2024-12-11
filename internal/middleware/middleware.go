package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
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
			// Change header to something like sessionid or something.
			log["something"] = params.Request.Header.Get("PIZZA")
			// time endresult is logged in milliseconds: 1000943 == 1 second.
			log["duration"] = time.Duration(params.Latency) / 1000

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func ValidateSecurityScheme(po out.PortUser, l *slog.Logger, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	clientToken, err := utils.GetBearerToken(l, input.RequestValidationInput.Request)
	if err != nil {
		l.Error(fmt.Sprintf("%v", err.Error()))
		return err
	}

	claims, err := utils.ValidateToken(l, clientToken)
	if err != nil {
		l.Debug(fmt.Sprintf("The '%v' has an incorrect token", claims.Username))
		myError := model.GetError("AUTH001")
		return errors.New(myError.Error)
	}

	if !slices.Contains(input.Scopes, claims.Role) {
		l.Debug(fmt.Sprintf("The '%v' is not port of the allowed roles/scopes.", claims.Role))
		myError := model.GetError("AUTH004")
		return errors.New(myError.Error)
	}

	user, _ := po.GetByName(claims.Username, ctx)
	if clientToken == user.Token {
		return nil
	} else {
		myError := model.GetError("AUTH002")
		return errors.New(myError.Error)
	}
}
