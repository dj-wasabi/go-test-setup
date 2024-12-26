package in

import (
	"context"

	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// TODO found in adapter/in/http/api/{organisations,authenticate,user}.go
type ApiUseCasesInterface interface {
	CreateOrganisation(ctx context.Context, c *model.Organisation) (*model.Organisation, *model.Error)
	GetAllOrganisations(ctx context.Context) (*model.ListOrganisations, *model.Error)
	AuthenticateLoginService(ctx context.Context, username, password, log_id string) (*model.AuthenticateToken, *model.Error)
	UserCreate(ctx context.Context, c *model.User) (*model.UserNoPassword, *model.Error)
	UserGet(ctx context.Context, userId, log_id string) (*model.UserNoPassword, *model.Error)
	// GetHealth(ctx context.Context)
}
