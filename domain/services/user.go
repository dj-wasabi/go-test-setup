package services

import (
	"context"
	"fmt"

	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (c *domainServices) UserCreate(ctx context.Context, command *model.User) *model.User {
	encryptPassword, _ := utils.HashPassword(&command.Password)
	command.Password = encryptPassword
	c.log.Info("The username", "key", command.Username)
	c.log.Info("This is an encrypted value", "key", encryptPassword)
	user := model.NewUser(command.Username, command.Password, command.Enabled, command.Roles)
	err := c.usr.Create(context.Background(), user)
	if err != nil {
		c.log.Error(fmt.Sprintf("%v", err))
	}
	return model.NewUser(user.GetUsername(), user.GetPassword(), user.GetEnabled(), user.GetRoles())
}
