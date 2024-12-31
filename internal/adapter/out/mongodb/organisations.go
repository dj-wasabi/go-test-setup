package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

type organisationService struct {
	// context    *context.Context
	logging    *slog.Logger
	repository *MongodbRepository
}

func NewOrganisationMongoRepo(con *mongo.Database, collection string) *MongodbRepository {
	return &MongodbRepository{
		DB:         con,
		Collection: con.Collection(collection),
	}
}

func NewOrganisationMongoService(repo *MongodbRepository, log *slog.Logger) out.PortOrganisationInterface {
	return &organisationService{
		logging:    log,
		repository: repo,
	}
}

// asasasasasa (These come from port/out/(interface))

func (mc *organisationService) CreateOrganisation(ctx context.Context, org *out.OrganizationPort) (*out.OrganizationPort, *model.Error) {
	org = out.NewOrganization(org.GetName(), org.GetDescription(), org.GetFqdn(), org.GetEnabled(), org.GetAdmins())

	if _, err := mc.repository.Collection.InsertOne(ctx, org); err != nil {
		var write_exc mongo.WriteException
		if !errors.As(err, &write_exc) {
			mc.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("%v", err))
			return nil, model.GetError("UNKNOWN", utils.GetLogId(ctx))
		}

		if write_exc.HasErrorCodeWithMessage(11000, "index: unique_name_fqdn_idx") {
			mc.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("User '%v' already exist, unique index violation.", org.Name))
			return nil, model.GetError("ORG0001", utils.GetLogId(ctx))
		}
	}
	return org, nil
}

func (mc *organisationService) GetAllOrganisations(ctx context.Context) ([]*out.OrganizationPort, *model.Error) {

	mc.logging.Debug("log_id", utils.GetLogId(ctx), "Get all available organisations from MongoDB")
	AllOrganisations := []*out.OrganizationPort{}

	cursor, err := mc.repository.Collection.Find(ctx, bson.D{})
	if err == mongo.ErrNoDocuments {
		mc.logging.Error("No document found")
	} else if err != nil {
		mc.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("Error in mongo %v", err.Error()))
		return nil, model.GetError("UNKNOWN", utils.GetLogId(ctx))
	}

	for cursor.Next(ctx) {
		result := new(out.OrganizationPort)
		err := cursor.Decode(&result)
		if err != nil {
			mc.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("Not able to decode database result with message: %v", err.Error()))
		}
		AllOrganisations = append(AllOrganisations, result)
	}

	return AllOrganisations, model.GetError("UNKNOWN", utils.GetLogId(ctx))

}
