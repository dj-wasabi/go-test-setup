package in

import (
	"time"
)

type ListOrganisationsIn struct {
	Organisations []OrganisationIn `json:"organisations"`
}

type OrganisationIn struct {
	Id          string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Fqdn        string    `json:"fqdn"`
	Admins      []string  `json:"admins,omitempty"`
}

type IOrganisation interface {
	GetId() string
	GetName() string
	GetDescription() string
	GetEnabled() bool
	GetCreated() time.Time
	GetUpdated() time.Time
	GetFqdn() string
	GetRoles() []string
}

func (o *OrganisationIn) GetId() string {
	return o.Id
}

func (o *OrganisationIn) GetName() string {
	return o.Name
}

func (o *OrganisationIn) GetDescription() string {
	return o.Description
}

func (o *OrganisationIn) GetEnabled() bool {
	return o.Enabled
}

func (o *OrganisationIn) GetCreated() time.Time {
	return o.CreatedAt
}

func (o *OrganisationIn) GetUpdated() time.Time {
	return o.UpdatedAt
}

func (o *OrganisationIn) GetFqdn() string {
	return o.Fqdn
}

func (o *OrganisationIn) GetAdmins() []string {
	return o.Admins
}

func NewOrganisationIn(name, description, fqdn string, enabled bool, admins []string) *OrganisationIn {
	return &OrganisationIn{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Fqdn:        fqdn,
		Admins:      admins,
	}
}
