package tokenstore

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (ts *tokenstoreService) Add(ctx context.Context, username, token string) error {
	ctx, span := tracer.Start(ctx, "Add")
	defer span.End()

	span.SetAttributes(
		attribute.String("tokenstore.file", "token"),
		attribute.String("tokenstore.function", "Add"),
		attribute.String("code.type", "adapter.out"),
	)

	ts.logging.Debug("log_id", utils.GetLogId(ctx), fmt.Sprintf("Set token data for username: %v", username))
	span.AddEvent("Add token to token store", trace.WithAttributes(
		attribute.String("username", username),
	))

	if err := ts.client.Set(ctx, username, token, 0).Err(); err != nil {
		ts.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("Error while adding token to tokenstore: %v", err.Error()))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (ts *tokenstoreService) Get(ctx context.Context, username string) (string, error) {
	ctx, span := tracer.Start(ctx, "Get")

	span.SetAttributes(
		attribute.String("tokenstore.file", "token"),
		attribute.String("tokenstore.function", "Get"),
		attribute.String("code.type", "adapter.out"),
	)

	ts.logging.Debug("log_id", utils.GetLogId(ctx), fmt.Sprintf("Get the token for username: %v", username))
	span.AddEvent("Get token from token store", trace.WithAttributes(
		attribute.String("username", username),
	))

	result, err := ts.client.Get(ctx, username).Result()
	if err != nil {
		ts.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("Error while getting token from tokenstore: %v", err.Error()))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	return result, nil
}
