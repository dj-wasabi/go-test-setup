package in

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// TODO found in adapter/in/http/api/{organisations,authenticate,user}.go
type ApiUseCases interface {
	CreateOrganisation(ctx context.Context, c *model.Organization) (*model.Organization, *model.Error)
	GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, *model.Error)
	AuthenticateLogin(ctx context.Context, username, password string) (*model.AuthenticationToken, *model.Error)
	UserCreate(ctx context.Context, c *model.User) (string, *model.Error)
}
