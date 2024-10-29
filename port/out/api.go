package out

import (
	"context"

	"werner-dijkerman.nl/test-setup/domain/model"
)

// TODO found in adapter/out/mongodb/{organisations,authenticate}.go
type OrganisationsDBPort interface {
	Create(ctx context.Context, org *OrganisationOutPort) (*OrganisationOutPort, error)
	GetAll(ctx context.Context) (*model.ListOrganisations, error)
	GetUserByName(username string, ctx context.Context) (*User, error)
}
