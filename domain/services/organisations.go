package services

import (
	"context"
	"log"

	"werner-dijkerman.nl/test-setup/domain/model"
)

func (c *domainServices) CreateOrganisation(ctx context.Context, command *model.Organization) *model.Organization {
	org := model.NewOrganization(command.Name, command.Description, command.Fqdn, command.Enabled, command.Admins)
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
	return allOrgs, err

}
