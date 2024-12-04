package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ListOrganisations struct {
	Organisations []Organisation `json:"organisations"`
}

var customOrganisationErrorMessages = map[string]string{
	"Username.required": "The field 'username' is required.",
	"Username.alphanum": "Only alphabetical and numerical characters are allowed.",
	"Username.min":      "The 'username' field needs a minimum amount of 6 characters.",
	"Username.max":      "The 'username' field has a maximum amount of 64 characters.",
	"Password.required": "The field 'password' is required.",
	"Password.min":      "The 'password' field needs a minimum amount of 6 characters.",
	"Password.max":      "The 'password' field has a maximum amount of 64 characters.",
}

type IOrganisation interface {
	GetId() string
	GetName() string
	GetDescription() string
	GetEnabled() bool
	GetFqdn() string
	GetTags() []string
}

func (o *Organisation) GetId() string {
	return o.Id
}

func (o *Organisation) GetName() string {
	return o.Name
}

func (o *Organisation) GetDescription() string {
	return o.Description
}

func (o *Organisation) GetEnabled() bool {
	return o.Enabled
}

func (o *Organisation) GetFqdn() string {
	return o.Fqdn
}

func (o *Organisation) GetAdmins() []string {
	return o.Admins
}

func NewOrganization(name, description, fqdn string, enabled bool, admins []string) (*Organisation, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := &Organisation{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Fqdn:        fqdn,
		Admins:      admins,
	}
	err := validate.Struct(e)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			field := fieldError.StructField()
			tag := fieldError.Tag()
			errorKey := fmt.Sprintf("%s.%s", field, tag)

			if message, keyFound := customOrganisationErrorMessages[errorKey]; keyFound {
				return &Organisation{}, errors.New(message)
			}
		}
	}
	return e, err
}
