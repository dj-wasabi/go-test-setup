package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

// adapter

type userService struct {
	logging    *slog.Logger
	repository *MongodbRepository
}

func NewUserMongoRepo(con *mongo.Database, collection string) *MongodbRepository {
	return &MongodbRepository{
		DB:         con,
		Collection: con.Collection(collection),
	}
}

func NewUserMongoService(repo *MongodbRepository, log *slog.Logger) out.PortUserInterface {
	return &userService{
		logging:    log,
		repository: repo,
	}
}

// Create an admin user if it doesn't exist while starting the application.
func NewAdminUser(pu out.PortUserInterface) error {
	ctx := context.Background()
	if _, error := pu.GetByName("administrator", ctx); error != nil {
		if error.Error == "no user found" {
			auth := utils.NewAuthentication()
			password := os.Getenv("INITIAL_ADMIN_PASSWORD")
			passwordHash, _ := auth.HashPassword(&password)
			adminUser := &out.UserPort{
				Username: "administrator",
				Password: passwordHash,
				Role:     "admin",
				Enabled:  true,
			}

			if _, createError := pu.Create(ctx, adminUser); createError != nil {
				return errors.New(createError.Error)
			}
			return nil
		}
	}
	return nil
}

func (uc *userService) Create(ctx context.Context, user *out.UserPort) (*out.UserPort, *model.Error) {
	ctx, span := tracer.Start(ctx, "UserCreate")
	defer span.End()

	span.SetAttributes(
		attribute.String("mongodb.file", "user"),
		attribute.String("mongodb.function", "Create"),
		attribute.String("code.type", "adapter.out"),
	)

	uc.logging.Debug(
		"log_id",
		utils.GetLogId(ctx),
		"Creating a new user and store it in MongoDB",
	)
	newUser := out.NewUser(user.GetUsername(), user.GetPassword(), user.GetRole(), user.GetEnabled(), string(user.GetOrgId()))

	span.AddEvent("Create user", trace.WithAttributes(
		attribute.String("username", user.Username),
		attribute.String("code.type", "adapter.out"),
	))
	uc.logging.Debug(
		"log_id",
		utils.GetLogId(ctx),
		fmt.Sprintf("Creating account with username '%v'", user.GetUsername()),
	)
	add, err := uc.repository.Collection.InsertOne(ctx, newUser)

	if err != nil {
		var write_exc mongo.WriteException
		if !errors.As(err, &write_exc) {
			uc.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("%v", err))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return &out.UserPort{}, model.GetError("UNKNOWN", utils.GetLogId(ctx))
		}

		if write_exc.HasErrorCodeWithMessage(11000, "index: unique_username_idx") {
			uc.logging.Error(
				"log_id",
				utils.GetLogId(ctx),
				fmt.Sprintf("User '%v' already exist, unique index violation.", user.GetUsername()),
			)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return &out.UserPort{}, model.GetError("USR0001", utils.GetLogId(ctx))
		}
	}

	userCreated := &out.UserPort{
		ID:        add.InsertedID.(primitive.ObjectID),
		Username:  newUser.Username,
		UpdatedAt: newUser.UpdatedAt,
		CreatedAt: newUser.CreatedAt,
		Enabled:   newUser.Enabled,
		Role:      newUser.Role,
		OrgId:     newUser.OrgId,
	}

	return userCreated, nil
}

func (uc *userService) GetByName(username string, ctx context.Context) (*out.UserPort, *model.Error) {
	ctx, span := tracer.Start(ctx, "UserGetByName")
	defer span.End()

	span.SetAttributes(
		attribute.String("mongodb.file", "user"),
		attribute.String("mongodb.function", "GetByName"),
		attribute.String("code.type", "adapter.out"),
	)
	uc.logging.Debug(
		"log_id",
		utils.GetLogId(ctx),
		fmt.Sprintf("Get user data from mongodb by looking for user with username: %v", username),
	)

	span.AddEvent("Finding user", trace.WithAttributes(
		attribute.String("username", username),
	))
	result := uc.repository.Collection.FindOne(ctx, bson.M{"username": username})
	if result.Err() == mongo.ErrNoDocuments {
		uc.logging.Info("log_id", utils.GetLogId(ctx), fmt.Sprintf("User '%v' not found.", username))
		return nil, model.NewError("no user found")
	}

	user := new(*out.UserPort)
	if err := result.Decode(&user); err != nil {
		uc.logging.Error(
			"log_id",
			utils.GetLogId(ctx),
			fmt.Sprintf("Error while decoding the user object, have error: '%v'", err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return *user, nil
}

func (uc *userService) GetById(userId string, ctx context.Context) (*out.UserPort, *model.Error) {
	ctx, span := tracer.Start(ctx, "UserGetById")
	defer span.End()

	span.SetAttributes(
		attribute.String("mongodb.file", "user"),
		attribute.String("mongodb.function", "GetById"),
	)
	uc.logging.Debug(
		"log_id",
		utils.GetLogId(ctx),
		fmt.Sprintf("Get user data from mongodb by looking for user with userid: %v", userId),
	)

	span.AddEvent("Finding user", trace.WithAttributes(
		attribute.String("userid", userId),
	))
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, model.GetError("USR0005", utils.GetLogId(ctx))
	}

	result := uc.repository.Collection.FindOne(ctx, bson.M{"_id": objectId})
	if result.Err() == mongo.ErrNoDocuments {
		uc.logging.Info("log_id", utils.GetLogId(ctx), fmt.Sprintf("User with id '%v' not found.", userId))
		return nil, model.NewError("no user found")
	}

	user := new(*out.UserPort)
	if err = result.Decode(&user); err != nil {
		uc.logging.Error(
			"log_id",
			utils.GetLogId(ctx),
			fmt.Sprintf("Error while decoding the user object, have error: '%v'", err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return *user, nil
}
