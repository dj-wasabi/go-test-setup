package model

import (
	"time"

	"werner-dijkerman.nl/test-setup/pkg/validator"
)

var customUserErrorMessages = map[string]string{
	"Username.required":         "The field 'username' is required.",
	"Username.alphanum":         "Only alphabetical and numerical characters are allowed.",
	"Username.min":              "The 'username' field needs a minimum amount of 6 characters.",
	"Username.max":              "The 'username' field has a maximum amount of 64 characters.",
	"Password.required":         "The field 'password' is required.",
	"Password.min":              "The 'password' field needs a minimum amount of 6 characters.",
	"Password.max":              "The 'password' field has a maximum amount of 64 characters.",
	"Password.validatePassword": "The password needs to have at least 1 uppercase, lowercase, number and any of the following characters: !\"#$%&'()*+,\\-./:;<=>?@[\\\\]^_`{|}~]",
	"Role.oneof":                "Only 'admin', 'write' or 'readonly' are allowed.",
	"OrgId.required":            "The field 'org_id' is required.",
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
	e := &User{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      role,
		OrgId:     orgid,
	}
	err := validator.CheckConfig(*e, customUserErrorMessages)
	return e, err
}
