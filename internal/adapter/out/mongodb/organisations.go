package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

// asasasasasa (These come from port/out/(interface))

func (mc *mongodbConnection) CreateOrganisation(ctx context.Context, org *out.OrganizationPort) (*out.OrganizationPort, *model.Error) {
	var mdbCollection string = "organisations"

	org = out.NewOrganization(org.GetName(), org.GetDescription(), org.GetFqdn(), org.GetEnabled(), org.GetAdmins())
	coll := mc.SetupCollection(mdbCollection)
	mc.Logging.Debug(fmt.Sprintf("We have the following name %v", org.Name))

	_, err := coll.InsertOne(context.TODO(), org)
	if err != nil {
		var write_exc mongo.WriteException
		if !errors.As(err, &write_exc) {
			mc.Logging.Error(fmt.Sprintf("%v", err))
			return nil, model.GetError("UNKNOWN")
		}

		if write_exc.HasErrorCodeWithMessage(11000, "index: unique_name_fqdn_idx") {
			mc.Logging.Error(fmt.Sprintf("User '%v' already exist, unique index violation.", org.Name))
			return nil, model.GetError("ORG0001")
		}
	}
	return org, nil
}

func (mc *mongodbConnection) GetAllOrganisations(ctx context.Context) ([]*out.OrganizationPort, *model.Error) {
	var mdbCollection string = "organisations"
	mc.Logging.Debug("Get all available organisations")

	coll := mc.SetupCollection(mdbCollection)
	AllOrganisations := []*out.OrganizationPort{}

	cursor, err := coll.Find(mc.Context, bson.D{})
	if err == mongo.ErrNoDocuments {
		mc.Logging.Error("No document found")
	} else if err != nil {
		mc.Logging.Error(fmt.Sprintf("Error in mongo %v", err.Error()))
		return nil, model.GetError("UNKNOWN")
	}

	for cursor.Next(mc.Context) {
		result := new(out.OrganizationPort)
		err := cursor.Decode(&result)
		if err != nil {
			mc.Logging.Error(fmt.Sprintf("Not able to decode database result with message: %v", err.Error()))
		}
		AllOrganisations = append(AllOrganisations, result)
	}

	return AllOrganisations, model.GetError("UNKNOWN")
}
