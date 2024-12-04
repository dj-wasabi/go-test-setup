package out

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

type UserPort struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Enabled   bool               `bson:"enabled"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Role      string             `bson:"role"`
	Token     string             `bson:"token"`
	OrgId     string             `bson:"organisation_id"`
}

type IUser interface {
	GetId() primitive.ObjectID
	GetOrgId() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetCreated() time.Time
	GetUpdated() time.Time
	GetRole() string
	GetToken() string
}

func (o *UserPort) GetId() primitive.ObjectID {
	return o.ID
}

func (o *UserPort) GetOrgId() string {
	return o.OrgId
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

func (o *UserPort) GetRole() string {
	return o.Role
}

func (o *UserPort) GetToken() string {
	return o.Token
}

func NewUser(username, password, role string, enabled bool, orgid string) *UserPort {
	return &UserPort{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      role,
		OrgId:     orgid,
	}
}

// TODO found in adapter/out/mongodb/{organisations,authenticate}.go
type PortUser interface {
	Create(ctx context.Context, user *UserPort) (*UserPort, *model.Error)
	GetByName(username string, ctx context.Context) (*UserPort, *model.Error)
	UpdateToken(ctx context.Context, token, username string) bool
}
