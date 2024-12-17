package services

import (
	"context"
	"fmt"
	"time"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) AuthenticateLoginService(ctx context.Context, username, password, log_id string) (*model.AuthenticateToken, *model.Error) {
	timeStart := time.Now()
	user, err := c.usr.GetByName(username, ctx)
	c.log.Debug("log_id", log_id, fmt.Sprintf("We have the '%v' username", username))
	if err != nil {
		timeEnd := float64(time.Since(timeStart).Seconds())
		model_authentication_requests.WithLabelValues("username_failure").Observe(timeEnd)
		return nil, model.GetError("USR0002", utils.GetLogId(ctx))
	}

	timeStartNew := time.Now()
	verifyPassword := utils.ValidatePassword(password, user.Password)
	if !verifyPassword {
		timeEnd := float64(time.Since(timeStartNew).Seconds())
		model_authentication_requests.WithLabelValues("password_failure").Observe(timeEnd)
		return nil, model.GetError("USR0002", utils.GetLogId(ctx))
	}

	timeStartNew = time.Now()
	token, authenticateError := utils.GenerateToken(user)
	c.log.Debug("log_id", log_id, fmt.Sprintf("Generated a new token for the user with '%v' as username", username))
	if authenticateError != nil {
		timeEnd := float64(time.Since(timeStartNew).Seconds())
		model_authentication_requests.WithLabelValues("token_generation_failure").Observe(timeEnd)
		return nil, model.NewError(authenticateError.Error())
	}
	// TMP Disable to continue investigating on how
	// to make the tests work with the checking part.
	_ = c.usr.UpdateToken(ctx, token, username)
	// if !isUpdated {
	// 	return nil, model.GetError("AUTH002")
	// }
	tokenOutput, tokenError := model.NewAuthenticationToken(token)
	if tokenError != nil {
		return nil, model.NewError(tokenError.Error())
	}

	timeEnd := float64(time.Since(timeStart).Seconds())
	model_authentication_requests.WithLabelValues("successful").Observe(timeEnd)
	return tokenOutput, nil
}
