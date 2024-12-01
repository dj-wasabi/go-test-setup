package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ListOrganisations struct {
	Organisations []Organisation `json:"organisations"`
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

func NewOrganization(name, description, fqdn string, enabled bool, admins []string) *Organisation {
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
		fmt.Println(err.Error())
		return &Organisation{}
	}
	return e
}
