package out

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

type OrganizationPort struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Enabled     bool               `bson:"enabled"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Fqdn        string             `bson:"fqdn"`
	Admins      []string           `bson:"admins"`
}

type IOrganisation interface {
	GetId() primitive.ObjectID
	GetName() string
	GetDescription() string
	GetEnabled() bool
	GetFqdn() string
	GetTags() []string
}

func (o *OrganizationPort) GetId() primitive.ObjectID {
	return o.Id
}

func (o *OrganizationPort) GetName() string {
	return o.Name
}

func (o *OrganizationPort) GetDescription() string {
	return o.Description
}

func (o *OrganizationPort) GetEnabled() bool {
	return o.Enabled
}

func (o *OrganizationPort) GetFqdn() string {
	return o.Fqdn
}

func (o *OrganizationPort) GetAdmins() []string {
	return o.Admins
}

func NewOrganization(name, description, fqdn string, enabled bool, admins []string) *OrganizationPort {
	return &OrganizationPort{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Fqdn:        fqdn,
		Admins:      admins,
	}
}

// TODO found in adapter/out/mongodb/{organisations,authenticate}.go
type PortOrganisation interface {
	CreateOrganisation(ctx context.Context, org *OrganizationPort) (*OrganizationPort, *model.Error)
	GetAllOrganisations(ctx context.Context) ([]*OrganizationPort, *model.Error)
}
