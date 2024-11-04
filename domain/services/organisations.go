package services

import (
	"context"
	"log"

	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/port/out"
)

func (c *domainServices) CreateOrganisation(ctx context.Context, command *model.Organization) *model.Organization {
	org := out.NewOrganization(command.Name, command.Description, command.Fqdn, command.Enabled, command.Admins)
	org, err := c.org.CreateOrganisation(context.Background(), org)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return model.NewOrganization(org.GetName(), org.GetDescription(), org.GetFqdn(), org.GetEnabled(), org.GetAdmins())
}

func (c *domainServices) GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, error) {
	allOrgs, err := c.org.GetAllOrganisations(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	AllOrganisations := &model.ListOrganisations{}
	for _, x := range allOrgs {
		result := new(model.Organization)
		result.ID = x.Id.Hex()
		result.Name = x.Name
		result.Description = x.Description
		result.Admins = x.Admins
		result.Fqdn = x.Fqdn
		result.CreatedAt = x.CreatedAt
		result.UpdatedAt = x.UpdatedAt
		AllOrganisations.Organisations = append(AllOrganisations.Organisations, *result)
	}

	return AllOrganisations, err
}
