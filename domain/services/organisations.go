package services

import (
	"context"
	"log"

	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/port/in"
	"werner-dijkerman.nl/test-setup/port/out"
)

func (c *DBServices) CreateOrganisation(ctx context.Context, command *in.OrganisationInPort) (*model.Organization, error) {
	org := out.NewOrganisationOutPort(command.Name, command.Description, command.Fqdn, command.Enabled, command.Admins)
	org, err := c.db.Create(context.Background(), org)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return model.NewOrganization(org.GetName(), org.GetDescription(), org.GetFqdn(), org.GetEnabled(), org.GetAdmins())
}

func (c *DBServices) GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, error) {
	allOrgs, err := c.db.GetAll(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return allOrgs, err

}
