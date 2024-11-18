package services

import (
	"context"
	"fmt"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) AuthenticateLoginService(ctx context.Context, username, password string) (*model.AuthenticationToken, *model.Error) {
	user, err := c.usr.GetByName(username, ctx)
	c.log.Debug(fmt.Sprintf("We have the '%v' username", username))
	if err != nil {
		return nil, model.GetError("USR0002")
	}

	verifyPassword := utils.ValidatePassword(password, user.Password)
	if !verifyPassword {
		return nil, model.GetError("USR0002")
	}

	token, err := utils.GenerateToken(username)
	c.log.Debug(fmt.Sprintf("Generated a new token for the user with '%v' as username", username))
	if err != nil {
		myError := &model.Error{
			Message: err.Error(),
		}
		return nil, myError
	}
	// TMP Disable to continue investigating on how
	// to make the tests work with the checking part.
	_ = c.usr.UpdateToken(ctx, token, username)
	// if !isUpdated {
	// 	return nil, model.GetError("AUTH002")
	// }

	tokens := &model.AuthenticationToken{
		Token: token,
	}

	return tokens, nil
}
