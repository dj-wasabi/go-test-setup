package services

import (
	"context"
	"fmt"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

func (c *domainServices) CreateOrganisation(ctx context.Context, command *model.Organization) (*model.Organization, *model.Error) {
	org := out.NewOrganization(command.Name, command.Description, command.Fqdn, command.Enabled, command.Admins)
	org, err := c.org.CreateOrganisation(context.Background(), org)
	if err != nil {
		return nil, err
	}

	newOrg := &model.Organization{
		Name:        org.Name,
		Description: org.Description,
		UpdatedAt:   org.UpdatedAt,
		CreatedAt:   org.CreatedAt,
		Admins:      org.Admins,
		Fqdn:        org.Fqdn,
	}

	return newOrg, nil
}

func (c *domainServices) GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, *model.Error) {
	allOrgs, err := c.org.GetAllOrganisations(context.Background())
	if err != nil {
		err = &model.Error{
			Message: fmt.Sprintf("%v", err),
		}
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
