package services

import (
	"context"

	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
	"werner-dijkerman.nl/test-setup/port/out"
)

func (c *domainServices) UserCreate(ctx context.Context, command *model.User) string {
	encryptPassword, _ := utils.HashPassword(&command.Password)
	command.Password = encryptPassword
	user := out.NewUser(command.Username, command.Password, command.Enabled, command.Roles)
	message := c.usr.Create(context.Background(), user)
	return message
}
