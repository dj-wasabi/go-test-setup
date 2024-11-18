package model

import "github.com/go-playground/validator/v10"

type AuthenticationRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type AuthenticationToken struct {
	Token string `json:"token" validate:"required"`
}

type Authentication interface {
	GetUsername() string
	GetPassword() string
}

func (a *AuthenticationRequest) GetUsername() string {
	return a.Username
}

func (a *AuthenticationRequest) GetPassword() string {
	return a.Password
}

func NewAuthenticationToken(token string) (*AuthenticationToken, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	t := &AuthenticationToken{
		Token: token,
	}

	err := validate.Struct(t)
	if err != nil {
		return &AuthenticationToken{}, err
	}
	return t, nil
}
