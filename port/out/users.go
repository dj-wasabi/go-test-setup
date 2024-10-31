package out

import (
	"context"

	"werner-dijkerman.nl/test-setup/domain/model"
)

// TODO found in adapter/out/mongodb/{organisations,authenticate}.go
type PortUser interface {
	Create(ctx context.Context, user *model.User) error
	GetByName(username string, ctx context.Context) (*model.User, error)
	UpdateToken(ctx context.Context, token, username string) bool
}
