package utils

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

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
		ts := &mockPortStoreInterface{
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
		ts := &mockPortStoreInterface{
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
		ts := &mockPortStoreInterface{
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

func TestValidatePassword(t *testing.T) {
	t.Run("ValidatePassWord", func(t *testing.T) {
		password := "dummy"
		newAuthentication := NewAuthentication()
		NewMockAuthentication := NewMockAuthentication()

		hashPassword, _ := NewMockAuthentication.HashPassword(&password)
		validateBool, validateError := newAuthentication.ValidatePassword(password, hashPassword)

		assert.Nil(t, validateError)
		assert.Equal(t, validateBool, true)
	})
}

func TestGenerateToken(t *testing.T) {
	mockUserPort := &mockUserPort{}
	user := &out.UserPort{
		Username: "test_user",
		OrgId:    "org1",
		Role:     "admin",
	}
	mockUserPort.On("GetById", user.ID).Return(user, nil)

	token, err := GenerateToken(user)
	if err != nil {
		t.Errorf("GenerateToken failed: %v", err)
	}

	// Validate the token
	claims := &AuthenticationDetails{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		t.Errorf("Invalid token: %v", err)
	}

	if claims.Username != user.Username {
		t.Errorf("Username in token does not match user: %s != %s", claims.Username, user.Username)
	}
	if claims.ExpiresAt != time.Now().Local().Add(time.Hour*time.Duration(24)).Unix() {
		t.Errorf("Expiry time in token is not correct: %d != %d", claims.ExpiresAt, time.Now().Local().Add(time.Hour*time.Duration(24)).Unix())
	}
}

// Mock stuff
type User struct {
	Enabled bool
}

type mockUserPort struct {
	users map[string]*User
	mock.Mock
}

type mockPortStoreInterface struct {
	tokens map[string]string
}

type mockAuthentication struct {
	mock.Mock
}

func NewMockAuthentication() mockAuthentication {
	return mockAuthentication{}
}

func (m *mockAuthentication) HashPassword(password *string) (string, error) {
	// Return encoded string for password: "dummy".
	return "$2a$14$Zdf1Ws4XFtyWOVO9HL2H/e8Jtz5pUlJFUbQa5TCQwTxq/r9iibqoG", nil
}

func (m *mockUserPort) GetByName(username string, ctx context.Context) (*out.UserPort, *model.Error) {
	if _, exists := m.users[username]; exists {
		return &out.UserPort{}, nil
	}
	return nil, model.GetError("UNKNOWN", "someid")
}

func (m *mockUserPort) GetById(string) (*out.UserPort, error) {
	args := m.Called(m)
	if res := args.Get(0); res != nil {
		return res.(*out.UserPort), nil
	}
	return nil, args.Error(1)
}

func (m *mockUserPort) Create(ctx context.Context, po *out.UserPort) (*out.UserPort, *model.Error) {
	username := po.Username
	if _, exists := m.users[username]; exists {
		return po, nil
	}
	return nil, model.GetError("UNKNOWN", "someid")
}

func (m *mockUserPort) UpdateToken(ctx context.Context, username, something string) bool {
	if _, exists := m.users[username]; exists {
		return false
	}
	return true
}

func (m *mockPortStoreInterface) Get(ctx context.Context, username string) (string, error) {
	if token, exists := m.tokens[username]; exists {
		return token, nil
	}
	return "", errors.New("token not found")
}

func (m *mockPortStoreInterface) Add(ctx context.Context, username, token string) error {
	if _, exists := m.tokens[username]; exists {
		return nil
	}
	return errors.New("token not found")
}
