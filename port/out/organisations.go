package out

import (
	"time"
)

type ListOrganisationOutPort struct {
	Organisations []OrganisationOutPort `json:"organisations"`
}

type OrganisationOutPort struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Fqdn        string    `json:"fqdn"`
	Admins      []string  `json:"admins,omitempty"`
}

type IOrganisationOutPort interface {
	GetName() string
	GetDescription() string
	GetEnabled() bool
	GetFqdn() string
	GetAdmins() []string
}

func (e *OrganisationOutPort) GetName() string {
	return e.Name
}

func (e *OrganisationOutPort) GetEnabled() bool {
	return e.Enabled
}

func (e *OrganisationOutPort) GetDescription() string {
	return e.Description
}

func (e *OrganisationOutPort) GetFqdn() string {
	return e.Fqdn
}

func (e *OrganisationOutPort) GetAdmins() []string {
	return e.Admins
}

func NewOrganisationOutPort(name, description, fqdn string, enabled bool, admins []string) *OrganisationOutPort {
	return &OrganisationOutPort{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Fqdn:        fqdn,
		Admins:      admins,
	}
}
