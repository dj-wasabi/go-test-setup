package in

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// TODO found in adapter/in/http/api/{organisations,authenticate,user}.go
type ApiUseCases interface {
	CreateOrganisation(ctx context.Context, c *model.Organisation) (*model.Organisation, *model.Error)
	GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, *model.Error)
	AuthenticateLoginService(ctx context.Context, username, password string) (*model.AuthenticatePostResponse, *model.Error)
	UserCreate(ctx context.Context, c *model.User) (string, *model.Error)
	// GetHealth(ctx context.Context)
}
