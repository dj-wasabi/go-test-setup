package services

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) UserCreate(ctx context.Context, command *model.User) (*model.UserNoPassword, *model.Error) {
	encryptPassword, _ := utils.HashPassword(&command.Password)
	command.Password = encryptPassword
	user := out.NewUser(command.Username, command.Password, command.Role, command.Enabled, command.OrgId)
	newUser, err := c.usr.Create(context.Background(), user)

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
