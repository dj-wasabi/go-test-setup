package out

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

type UserPort struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Enabled   bool               `bson:"enabled"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Roles     []string           `bson:"roles"`
}

type IUser interface {
	GetId() primitive.ObjectID
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetCreated() time.Time
	GetUpdated() time.Time
	GetRoles() []string
}

func (o *UserPort) GetId() primitive.ObjectID {
	return o.ID
}

func (o *UserPort) GetUsername() string {
	return o.Username
}

func (o *UserPort) GetPassword() string {
	return o.Password
}

func (o *UserPort) GetEnabled() bool {
	return o.Enabled
}

func (o *UserPort) GetCreated() time.Time {
	return o.CreatedAt
}

func (o *UserPort) GetUpdated() time.Time {
	return o.UpdatedAt
}

func (o *UserPort) GetRoles() []string {
	return o.Roles
}

func NewUser(username, password string, enabled bool, roles []string) *UserPort {
	return &UserPort{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Roles:     roles,
	}
}

// TODO found in adapter/out/mongodb/{organisations,authenticate}.go
type PortUser interface {
	Create(ctx context.Context, user *UserPort) (string, *model.Error)
	GetByName(username string, ctx context.Context) (*UserPort, error)
	UpdateToken(ctx context.Context, token, username string) bool
}
