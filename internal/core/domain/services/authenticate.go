package services

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) AuthenticateLoginService(ctx context.Context, username, password, log_id string) (*model.AuthenticateToken, *model.Error) {
	span := trace.SpanFromContext(ctx)
	timeStart := time.Now()
	user, err := c.usr.GetByName(username, ctx)
	c.log.Debug("log_id", log_id, fmt.Sprintf("We have the '%v' username", username))
	if err != nil {
		timeEnd := float64(time.Since(timeStart).Seconds())
		model_authentication_requests.WithLabelValues("username_failure").Observe(timeEnd)
		return nil, model.GetError("USR0002", utils.GetLogId(ctx))
	}

	timeStartNew := time.Now()
	auth := utils.NewAuthentication()

	span.AddEvent("Verifying password", trace.WithAttributes(attribute.String("code.type", "service")))
	if verifyPassword, _ := auth.ValidatePassword(password, user.Password); !verifyPassword {
		timeEnd := float64(time.Since(timeStartNew).Seconds())
		model_authentication_requests.WithLabelValues("password_validation_failure").Observe(timeEnd)
		return nil, model.GetError("USR0002", utils.GetLogId(ctx))
	}

	timeStartNew = time.Now()
	token, authenticateError := utils.GenerateToken(user)
	span.AddEvent("Generating a token", trace.WithAttributes(attribute.String("username", username), attribute.String("code.type", "service")))
	c.log.Debug("log_id", log_id, fmt.Sprintf("Generated a new token for the user with '%v' as username", username))

	if authenticateError != nil {
		timeEnd := float64(time.Since(timeStartNew).Seconds())
		model_authentication_requests.WithLabelValues("token_generation_failure").Observe(timeEnd)
		return nil, model.NewError(authenticateError.Error())
	}

	// Update tokenstore
	span.AddEvent("Adding token to tokenstore", trace.WithAttributes(attribute.String("code.type", "service")))
	if addError := c.token.Add(ctx, username, token); addError != nil {
		return nil, model.NewError(addError.Error())
	}

	tokenOutput, tokenError := model.NewAuthenticationToken(token)
	if tokenError != nil {
		return nil, model.NewError(tokenError.Error())
	}

	timeEnd := float64(time.Since(timeStart).Seconds())
	model_authentication_requests.WithLabelValues("successful").Observe(timeEnd)
	return tokenOutput, nil
}
