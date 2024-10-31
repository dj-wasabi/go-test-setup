package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"-" bson:"password"`
	Enabled   bool      `json:"enabled" bson:"enabled"`
	CreatedAt time.Time `json:"-" bson:"created_at"`
	UpdatedAt time.Time `json:"-" bson:"updated_at"`
	Roles     []string  `json:"roles,omitempty" bson:"roles"`
}

type IUser interface {
	GetId() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetRoles() []string
}

func (o *User) GetId() string {
	return o.ID
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

func (o *User) GetRoles() []string {
	return o.Roles
}

func NewUser(username, password string, enabled bool, roles []string) *User {
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := &User{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Roles:     roles,
	}
	err := validate.Struct(e)
	if err != nil {
		fmt.Println(err.Error())
		return &User{}
	}
	return e
}
