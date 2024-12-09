package model

import (
	"werner-dijkerman.nl/test-setup/pkg/validator"
)

var customAuthErrorMessages = map[string]string{
	"Token.required": "The field 'username' is required.",
}

type Authentication interface {
	GetUsername() string
	GetPassword() string
}

func (a *AuthenticateRequest) GetUsername() string {
	return a.Username
}

func (a *AuthenticateRequest) GetPassword() string {
	return a.Password
}

func ValidateAuthenticationData(at *AuthenticateRequest) error {
	err := validator.CheckConfig(*at, customUserErrorMessages)
	return err
}

func NewAuthenticationToken(token string) (*AuthenticateToken, error) {
	t := &AuthenticateToken{
		Token: token,
	}

	err := validator.CheckConfig(*t, customAuthErrorMessages)
	return t, err
}
