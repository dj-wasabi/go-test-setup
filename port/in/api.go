package in

import (
	"context"

	"werner-dijkerman.nl/test-setup/domain/model"
)

// TODO found in adapter/in/http/api/{organisations,}.go
type ApiUseCases interface {
	CreateOrganisation(ctx context.Context, c *OrganisationInPort) (*model.Organization, error)
	GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, error)
	AuthenticateLogin(ctx context.Context, username, password string)
}
