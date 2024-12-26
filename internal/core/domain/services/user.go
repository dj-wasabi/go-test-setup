package services

import (
	"context"
	"fmt"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) UserCreate(ctx context.Context, command *model.User) (*model.UserNoPassword, *model.Error) {
	c.log.Debug("log_id", utils.GetLogId(ctx), "Domain logic to create a new user")

	auth := utils.NewAuthentication()
	encryptPassword, _ := auth.HashPassword(&command.Password)
	command.Password = encryptPassword

	user := out.NewUser(command.Username, command.Password, command.Role, command.Enabled, command.OrgId)
	newUser, err := c.usr.Create(ctx, user)

	returnUser := &model.UserNoPassword{
		Id:        newUser.ID.Hex(),
		Username:  newUser.Username,
		UpdatedAt: newUser.UpdatedAt,
		CreatedAt: newUser.CreatedAt,
		Enabled:   newUser.Enabled,
		Role:      newUser.Role,
		OrgId:     newUser.OrgId,
	}

	return returnUser, err
}

func (c *domainServices) UserGet(ctx context.Context, userId, log_id string) (*model.UserNoPassword, *model.Error) {

	user, err := c.usr.GetById(userId, ctx)
	c.log.Debug("log_id", log_id, fmt.Sprintf("We have the '%v' userId", userId))
	if err != nil {
		c.log.Debug("log_id", log_id, fmt.Sprintf("User with id '%v' not found in the database", userId))
		return nil, model.GetError("USR0004", utils.GetLogId(ctx))
	}

	returnUser := &model.UserNoPassword{
		Id:        user.ID.Hex(),
		Username:  user.Username,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
		Enabled:   user.Enabled,
		Role:      user.Role,
		OrgId:     user.OrgId,
	}

	return returnUser, nil
}
