package model

import "github.com/go-playground/validator/v10"

type Authentication interface {
	GetUsername() string
	GetPassword() string
}

func (a *AuthenticatePostRequest) GetUsername() string {
	return a.Username
}

func (a *AuthenticatePostRequest) GetPassword() string {
	return a.Password
}

func NewAuthenticationToken(token string) (*AuthenticateToken, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	t := &AuthenticateToken{
		Token: token,
	}

	err := validate.Struct(t)
	if err != nil {
		return &AuthenticateToken{}, err
	}
	return t, nil
}
