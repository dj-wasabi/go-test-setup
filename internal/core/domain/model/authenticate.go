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

func NewAuthenticationToken(token string) (*AuthenticatePostResponse, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	t := &AuthenticatePostResponse{
		Token: token,
	}

	err := validate.Struct(t)
	if err != nil {
		return &AuthenticatePostResponse{}, err
	}
	return t, nil
}
