package tokenstore

import (
	"context"
	"fmt"

	"werner-dijkerman.nl/test-setup/pkg/utils"
)

func (ts *tokenstoreService) Add(ctx context.Context, username, token string) error {
	ts.logging.Debug("log_id", utils.GetLogId(ctx), fmt.Sprintf("Set token data for username: %v", username))

	if err := ts.client.Set(ctx, username, token, 0).Err(); err != nil {
		ts.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("Error while adding token to tokenstore: %v", err.Error()))
		return err
	}
	return nil
}

func (ts *tokenstoreService) Get(ctx context.Context, username string) (string, error) {
	ts.logging.Debug("log_id", utils.GetLogId(ctx), fmt.Sprintf("Get the token for username: %v", username))

	result, err := ts.client.Get(ctx, username).Result()
	if err != nil {
		ts.logging.Error("log_id", utils.GetLogId(ctx), fmt.Sprintf("Error while getting token from tokenstore: %v", err.Error()))
		return "", err
	}
	return result, nil
}
