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
	"github.com/google/uuid"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

var logid string = "X-APP-LOG-ID"

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
			log["logId"] = params.Request.Header.Get("X-APP-LOG-ID")
			// time endresult is logged in milliseconds: 1000943 == 1 second.
			log["duration"] = time.Duration(params.Latency) / 1000

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func ValidateSecurityScheme(po out.PortUser, l *slog.Logger, input *openapi3filter.AuthenticationInput) error {
	logId := input.RequestValidationInput.Request.Header.Get(logid)
	if logId == "" {
		logId_string := uuid.New()
		logId = logId_string.String()
		input.RequestValidationInput.Request.Header.Set(logid, logId_string.String())
	}
	ctx := utils.NewContextWrapper(context.TODO(), logId).Build()

	ad, clientToken, err := utils.GetAuthenticationDetails(l, input.RequestValidationInput.Request, logId)
	if err != nil {
		l.Error("log_id", logId, fmt.Sprintf("%v", err.Error()))
		return err
	}

	err = ad.Validate(l, logId)
	if err != nil {
		myError := model.GetError("AUTH001", logId)
		return errors.New(myError.Error)
	}

	if !slices.Contains(input.Scopes, ad.GetRole()) {
		l.Debug("log_id", logId, fmt.Sprintf("The '%v' is not port of the allowed roles/scopes.", ad.GetRole()))
		myError := model.GetError("AUTH004", logId)
		return errors.New(myError.Error)
	}

	user, _ := po.GetByName(ad.GetUsername(), ctx)
	if clientToken == user.Token {
		l.Debug("log_id", logId, fmt.Sprintf("Successfully validated token for '%v'", ad.GetUsername()))
		return nil
	} else {
		myError := model.GetError("AUTH002", logId)
		return errors.New(myError.Error)
	}
}
