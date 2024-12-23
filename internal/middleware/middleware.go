package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

var logHeaderString string = "X-APP-LOG-ID"

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

func ensureLogID(input *openapi3filter.AuthenticationInput) string {
	logID := input.RequestValidationInput.Request.Header.Get(logHeaderString)
	if logID == "" {
		newLogID := uuid.New().String()
		input.RequestValidationInput.Request.Header.Set(logHeaderString, newLogID)
		logID = newLogID
	}
	return logID
}

func ValidateSecurityScheme(po out.PortUser, ts out.PortStore, l *slog.Logger, input *openapi3filter.AuthenticationInput) error {
	logId := ensureLogID(input)
	ctx := utils.NewContextWrapper(context.TODO(), logId).Build()

	// Get token information and validate the token
	adToken, providedToken, err := utils.GetAuthenticationDetails(l, input.RequestValidationInput.Request, logId)
	if err != nil {
		l.Error("log_id", logId, fmt.Sprintf("Failed to get authentication details: %v", err.Error()))
		return err
	}

	if validateError := adToken.Validate(l, logId); validateError != nil {
		return utils.HandleAuthError("AUTH001", logId, l)
	}

	if !slices.Contains(input.Scopes, adToken.GetRole()) {
		l.Debug("log_id", logId, fmt.Sprintf("The '%v' is not port of the allowed roles/scopes.", adToken.GetRole()))
		return utils.HandleAuthError("AUTH004", logId, l)
	}

	// Verify user status
	if err := utils.ValidateUserStatus(po, ctx, adToken.GetUsername(), logId, l); err != nil {
		return err
	}

	// Validate token against stored token
	if err := utils.ValidateStoredToken(ts, ctx, adToken.GetUsername(), providedToken, logId, l); err != nil {
		return err
	}

	l.Debug("log_id", logId, fmt.Sprintf("Successfully validated token for '%v'", adToken.GetUsername()))
	return nil
}
