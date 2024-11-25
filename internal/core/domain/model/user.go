package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id        string
	Username  string
	Password  string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string
}

type IUser interface {
	GetId() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetRole() string
}

func (o *User) GetId() string {
	return o.Id
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

func NewUser(username, password, role string, enabled bool) (*User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := &User{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      role,
	}
	err := validate.Struct(e)
	if err != nil {
		fmt.Println(err.Error())
		return &User{}, err
	}
	return e, err
}
