package services

import (
	"context"

	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) AuthenticateLogin(ctx context.Context, username, password string) (*model.AuthenticationToken, error) {
	user, err := c.usr.GetByName(username, ctx)
	if err != nil {
		return nil, err
	}

	// result := new(*out.AuthenticationToken)

	verifyPassword := utils.ValidatePassword(password, user.Password)
	if !verifyPassword {
		return nil, err
	}

	token, err := utils.GenerateToken(username)
	if err != nil {
		return nil, err
	}
	isUpdated := c.usr.UpdateToken(ctx, token, username)
	if !isUpdated {
		return nil, err
	}
	// token, err := out.NewAuthenticationToken(pizza)

	tokens := &model.AuthenticationToken{
		Token: token,
	}

	return tokens, nil
}
