package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
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

// Custom check to validate the provided password. Could not find an easy way to rely
// on the OpenAPI/Struct Validator and making our own validator would be the best way.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasLower = regexp.MustCompile(`[[:lower:]]`)
	var hasUpper = regexp.MustCompile(`[[:upper:]]`)
	var hasNumber = regexp.MustCompile(`[[:digit:]]`)
	var hasCharacters = regexp.MustCompile(`[[:graph:]]`)

	if !hasLower.MatchString(password) || !hasUpper.MatchString(password) || !hasNumber.MatchString(password) || !hasCharacters.MatchString(password) {
		return false
	}

	return true
}

func NewUser(username, password, role string, enabled bool, orgid string) (*User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("validatePassword", validatePassword)

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
