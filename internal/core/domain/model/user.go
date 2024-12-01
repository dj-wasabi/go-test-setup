package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type IUser interface {
	GetId() string
	GetOrgId() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetRole() string
}

func (o *User) GetId() string {
	return o.Id
}

func (o *User) GetOrgId() string {
	return o.OrgId
}

func (o *User) GetUsername() string {
	return o.Username
}

func (o *User) GetPassword() string {
	return o.Password
}

func (o *User) GetEnabled() bool {
	return o.Enabled
}

func (o *User) GetRole() string {
	return o.Role
}

func NewUser(username, password, role string, enabled bool, orgid string) (*User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := &User{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      role,
		OrgId:     orgid,
	}
	err := validate.Struct(e)
	if err != nil {
		fmt.Println(err.Error())
		return &User{}, err
	}
	return e, err
}
