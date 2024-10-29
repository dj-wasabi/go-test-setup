package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/port/out"
)

// asasasasasa (These come from port/out/(interface))
func (mc *mongodbConnection) Create(ctx context.Context, org *out.OrganisationOutPort) (*out.OrganisationOutPort, error) {
	var mdbCollection string = "organisation"
	mc.Logging.Debug(fmt.Sprintf("About to Create Organisations %v", org))

	name := org.GetName()
	org = out.NewOrganisationOutPort(org.GetName(), org.GetDescription(), org.GetFqdn(), org.GetEnabled(), org.GetAdmins())
	coll := mc.SetupCollection(mdbCollection)

	result := coll.FindOne(mc.Context, bson.M{"name": name})
	mc.Logging.Info(fmt.Sprintf("%v", result.Err()))
	if result.Err() == mongo.ErrNoDocuments {
		mc.Logging.Info(org.Name)
		_, _ = coll.InsertOne(context.TODO(), org)
	}

	// defer mc.CloseConnection(mc.Cancel)

	return org, nil
}

func (mc *mongodbConnection) GetAll(ctx context.Context) (*model.ListOrganisations, error) {
	var mdbCollection string = "organisation"
	mc.Logging.Debug("Get all available organisations")

	coll := mc.SetupCollection(mdbCollection)
	AllOrganisations := &model.ListOrganisations{}

	filter := bson.D{}
	cursor, err := coll.Find(mc.Context, filter)

	if err == mongo.ErrNoDocuments {
		mc.Logging.Error("No document found")
	} else if err != nil {
		mc.Logging.Error(fmt.Sprintf("Error in mongo %v", err.Error()))
	}

	// defer cursor.Close(mc.Context)

	for cursor.Next(mc.Context) {
		result := new(model.Organization)
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		AllOrganisations.Organisations = append(AllOrganisations.Organisations, *result)
	}

	return AllOrganisations, err
}
