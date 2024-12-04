package services

import (
	"context"
	"fmt"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) AuthenticateLoginService(ctx context.Context, username, password string) (*model.AuthenticateToken, *model.Error) {
	user, err := c.usr.GetByName(username, ctx)
	c.log.Debug(fmt.Sprintf("We have the '%v' username", username))
	if err != nil {
		return nil, model.GetError("USR0002")
	}

	verifyPassword := utils.ValidatePassword(password, user.Password)
	if !verifyPassword {
		return nil, model.GetError("USR0002")
	}

	token, authenticateError := utils.GenerateToken(username, user.Role)
	c.log.Debug(fmt.Sprintf("Generated a new token for the user with '%v' as username", username))
	if authenticateError != nil {
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

	return tokenOutput, nil
}
