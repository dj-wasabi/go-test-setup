package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ListOrganisations struct {
	Organisations []Organization `json:"organisations"`
}

type Organization struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Fqdn        string    `json:"fqdn"`
	Admins      []string  `json:"admins,omitempty"`
}

type IOrganisation interface {
	GetId() string
	GetName() string
	GetDescription() string
	GetEnabled() bool
	GetFqdn() string
	GetTags() []string
}

func (o *Organization) GetId() string {
	return o.ID
}

func (o *Organization) GetName() string {
	return o.Name
}

func (o *Organization) GetDescription() string {
	return o.Description
}

func (o *Organization) GetEnabled() bool {
	return o.Enabled
}

func (o *Organization) GetFqdn() string {
	return o.Fqdn
}

func (o *Organization) GetAdmins() []string {
	return o.Admins
}

func NewOrganization(name, description, fqdn string, enabled bool, admins []string) *Organization {
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := &Organization{
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
		fmt.Println(err.Error())
		return &Organization{}
	}
	return e
}
