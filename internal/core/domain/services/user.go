package services

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) UserCreate(ctx context.Context, command *model.User) (*model.UserNoPassword, *model.Error) {
	c.log.Debug("log_id", utils.GetLogId(ctx), "Domain logic to create a new user")
	encryptPassword, _ := utils.HashPassword(&command.Password)
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
