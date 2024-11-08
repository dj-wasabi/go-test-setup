package services

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) AuthenticateLogin(ctx context.Context, username, password string) (*model.AuthenticationToken, *model.Error) {
	user, err := c.usr.GetByName(username, ctx)
	if err != nil {
		return nil, model.GetError("USR0002")
	}

	verifyPassword := utils.ValidatePassword(password, user.Password)
	if !verifyPassword {
		return nil, model.GetError("USR0002")
	}

	token, err := utils.GenerateToken(username)
	if err != nil {
		myError := &model.Error{
			Message: err.Error(),
		}
		return nil, myError
	}
	isUpdated := c.usr.UpdateToken(ctx, token, username)
	if !isUpdated {
		return nil, model.GetError("AUTH002")
	}

	tokens := &model.AuthenticationToken{
		Token: token,
	}

	return tokens, nil
}
