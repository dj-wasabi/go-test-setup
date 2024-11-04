package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/port/out"
)

// asasasasasa (These come from port/out/(interface))
func (mc *mongodbConnection) CreateOrganisation(ctx context.Context, org *out.OrganizationPort) (*out.OrganizationPort, error) {
	var mdbCollection string = "organisation"
	mc.Logging.Debug(fmt.Sprintf("About to Create Organisations %v", org))

	name := org.GetName()
	org = out.NewOrganization(org.GetName(), org.GetDescription(), org.GetFqdn(), org.GetEnabled(), org.GetAdmins())
	coll := mc.SetupCollection(mdbCollection)

	result := coll.FindOne(mc.Context, bson.M{"name": name})
	mc.Logging.Info(fmt.Sprintf("%v", result.Err()))
	if result.Err() == mongo.ErrNoDocuments {
		mc.Logging.Info(org.Name)
		_, _ = coll.InsertOne(context.TODO(), org)
	}

	return org, nil
}

func (mc *mongodbConnection) GetAllOrganisations(ctx context.Context) ([]*out.OrganizationPort, error) {
	var mdbCollection string = "organisation"
	mc.Logging.Debug("Get all available organisations")

	coll := mc.SetupCollection(mdbCollection)
	AllOrganisations := []*out.OrganizationPort{}

	cursor, err := coll.Find(mc.Context, bson.D{})
	if err == mongo.ErrNoDocuments {
		mc.Logging.Error("No document found")
	} else if err != nil {
		mc.Logging.Error(fmt.Sprintf("Error in mongo %v", err.Error()))
	}

	for cursor.Next(mc.Context) {
		result := new(out.OrganizationPort)
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		AllOrganisations = append(AllOrganisations, result)
	}

	return AllOrganisations, err
}
