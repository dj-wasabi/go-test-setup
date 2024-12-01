package in

type IOrganisation interface {
	GetName() string
	GetDescription() string
	GetEnabled() bool
	GetFqdn() string
	GetRoles() []string
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
		Fqdn:        fqdn,
		Admins:      admins,
	}
}
