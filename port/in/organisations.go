package in

import (
	"time"
)

type ListOrganisationInPort struct {
	Organisations OrganisationInPort
}

type OrganisationInPort struct {
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Fqdn        string
	Admins      []string
}

func NewListOrganisationInPort(org *OrganisationInPort) *ListOrganisationInPort {
	return &ListOrganisationInPort{Organisations: *org}
}

func NewOrganisationInPort(name, description, fqdn string, enabled bool, admins []string) *OrganisationInPort {
	return &OrganisationInPort{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		Fqdn:        fqdn,
		Admins:      admins,
	}
}
