package out

import (
	"context"

	"werner-dijkerman.nl/test-setup/domain/model"
)

// TODO found in adapter/out/mongodb/{organisations,authenticate}.go
type PortOrganisation interface {
	CreateOrganisation(ctx context.Context, org *model.Organization) (*model.Organization, error)
	GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, error)
}
