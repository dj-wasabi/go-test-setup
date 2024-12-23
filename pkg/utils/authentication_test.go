package utils

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

// var po out.PortUser

func TestHandleAuthError(t *testing.T) {
	t.Run("Returns error and logs it", func(t *testing.T) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
		logId := "test-log-id"
		err := HandleAuthError("AUTH001", logId, logger)

		assert.NotNil(t, err)
		assert.Equal(t, "Error while validating the token.", err.Error())
	})
}

func TestValidateStoredToken(t *testing.T) {
	t.Run("Tokens match", func(t *testing.T) {
		ts := &mockPortStore{
			tokens: map[string]string{
				"testuser": "valid-token",
			},
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
		ctx := context.Background()
		logId := "test-log-id"

		err := ValidateStoredToken(ts, ctx, "testuser", "valid-token", logId, logger)

		assert.Nil(t, err)
	})

	t.Run("Tokens do not match", func(t *testing.T) {
		ts := &mockPortStore{
			tokens: map[string]string{
				"testuser": "valid-token",
			},
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
		ctx := context.Background()
		logId := "test-log-id"

		err := ValidateStoredToken(ts, ctx, "testuser", "invalid-token", logId, logger)

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid token/user combination.", err.Error())
	})

	t.Run("Token not found", func(t *testing.T) {
		ts := &mockPortStore{
			tokens: map[string]string{},
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
		ctx := context.Background()
		logId := "test-log-id"

		err := ValidateStoredToken(ts, ctx, "unknownuser", "some-token", logId, logger)

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid token/user combination.", err.Error())
	})
}

// Mock implementations

type User struct {
	Enabled bool
}

type mockPortUser struct {
	users map[string]*User
}

func (m *mockPortUser) GetByName(username string, ctx context.Context) (*out.UserPort, *model.Error) {
	if _, exists := m.users[username]; exists {
		return &out.UserPort{}, nil
	}
	return nil, model.GetError("UNKNOWN", "someid")
}

func (m *mockPortUser) GetById(username string, ctx context.Context) (*out.UserPort, *model.Error) {
	if _, exists := m.users[username]; exists {
		return &out.UserPort{}, nil
	}
	return nil, model.GetError("UNKNOWN", "someid")
}

func (m *mockPortUser) Create(ctx context.Context, po *out.UserPort) (*out.UserPort, *model.Error) {
	username := po.Username
	if _, exists := m.users[username]; exists {
		return po, nil
	}
	return nil, model.GetError("UNKNOWN", "someid")
}

func (m *mockPortUser) UpdateToken(ctx context.Context, username, something string) bool {
	if _, exists := m.users[username]; exists {
		return false
	}
	return true
}

type mockPortStore struct {
	tokens map[string]string
}

func (m *mockPortStore) Get(ctx context.Context, username string) (string, error) {
	if token, exists := m.tokens[username]; exists {
		return token, nil
	}
	return "", errors.New("token not found")
}

func (m *mockPortStore) Add(ctx context.Context, username, token string) error {
	if _, exists := m.tokens[username]; exists {
		return nil
	}
	return errors.New("token not found")
}
