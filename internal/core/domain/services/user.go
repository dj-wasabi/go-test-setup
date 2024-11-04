package services

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) UserCreate(ctx context.Context, command *model.User) string {
	encryptPassword, _ := utils.HashPassword(&command.Password)
	command.Password = encryptPassword
	user := out.NewUser(command.Username, command.Password, command.Enabled, command.Roles)
	message := c.usr.Create(context.Background(), user)
	return message
}
