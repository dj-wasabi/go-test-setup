package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var customUserErrorMessages = map[string]string{
	"Username.required": "The field 'username' is required.",
	"Username.alphanum": "Only alphabetical and numerical characters are allowed.",
	"Username.min":      "The 'username' field needs a minimum amount of 6 characters.",
	"Username.max":      "The 'username' field has a maximum amount of 64 characters.",
	"Password.required": "The field 'password' is required.",
	"Password.min":      "The 'password' field needs a minimum amount of 6 characters.",
	"Password.max":      "The 'password' field has a maximum amount of 64 characters.",
}

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
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			field := fieldError.StructField()
			tag := fieldError.Tag()
			errorKey := fmt.Sprintf("%s.%s", field, tag)

			if message, keyFound := customUserErrorMessages[errorKey]; keyFound {
				return &User{}, errors.New(message)
			}
		}
	}
	return e, err
}
