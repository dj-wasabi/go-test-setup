package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
)

type AuthenticationDetails struct {
	Username string
	UserId   string
	OrgId    string
	Role     string
	jwt.StandardClaims
}

type AuthenticationImpl struct {
	Authentication
}

type Authentication interface {
	ValidatePassword(string, string) (bool, error)
	HashPassword(*string) (string, error)
	GetAuthenticationDetails(*http.Request, string) (*AuthenticationDetails, string, error)
}

var (
	ErrNoAuthHeader             = errors.New("authorization header is missing")
	ErrInvalidAuthHeader        = errors.New("authorization header is malformed")
	ErrInvalidToken             = errors.New("token is invalid")
	ErrTokenExpired             = errors.New("the token is expired")
	SECRET_KEY           string = os.Getenv("SECRET_KEY")
)

func (ad *AuthenticationDetails) GetUsername() string {
	return ad.Username
}

func (ad *AuthenticationDetails) GetUserId() string {
	return ad.UserId
}

func (ad *AuthenticationDetails) GetOrgId() string {
	return ad.OrgId
}

func (ad *AuthenticationDetails) GetRole() string {
	return ad.Role
}

func (ad *AuthenticationDetails) GetToken() jwt.Claims {
	return ad.StandardClaims
}

func NewAuthentication() Authentication {
	return AuthenticationImpl{}
}

func (ai AuthenticationImpl) GetAuthenticationDetails(req *http.Request, logId string) (*AuthenticationDetails, string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return &AuthenticationDetails{}, "", ErrNoAuthHeader
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return &AuthenticationDetails{}, "", ErrInvalidAuthHeader
	}

	myTokenString := strings.TrimPrefix(authHdr, prefix)
	jwtToken, err := jwt.ParseWithClaims(myTokenString, &AuthenticationDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, "", err
	}
	newClaims, ok := jwtToken.Claims.(*AuthenticationDetails)
	if !ok {
		return nil, "", ErrInvalidToken
	}

	claim := &AuthenticationDetails{
		Username:       newClaims.Username,
		UserId:         newClaims.UserId,
		OrgId:          newClaims.OrgId,
		Role:           newClaims.Role,
		StandardClaims: newClaims.StandardClaims,
	}

	return claim, myTokenString, nil
}

func (ad *AuthenticationDetails) Validate(l *slog.Logger, logId string) error {
	if ad.ExpiresAt < time.Now().Local().Unix() {
		return ErrTokenExpired
	}
	return nil
}

func (ai AuthenticationImpl) ValidatePassword(providedpassword, storedpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(storedpassword), []byte(providedpassword))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ai AuthenticationImpl) HashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GenerateToken(user *out.UserPort) (signedToken string, err error) {
	claims := &AuthenticationDetails{
		Username: user.Username,
		UserId:   user.ID.Hex(),
		OrgId:    user.OrgId,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	return token, err
}

func ValidateUserStatus(po out.PortUserInterface, ctx context.Context, username, logId string, l *slog.Logger) error {
	user, err := po.GetByName(username, ctx)
	if err != nil {
		l.Error("log_id", logId, fmt.Sprintf("Failed to fetch user '%v': %v", username, err))
		return HandleAuthError("AUTH002", logId, l)
	}

	if !user.Enabled {
		l.Debug("log_id", logId, fmt.Sprintf("User '%v' is not enabled.", username))
		return HandleAuthError("AUTH005", logId, l)
	}
	return nil
}

func ValidateStoredToken(ts out.PortStoreInterface, ctx context.Context, username, providedToken, logId string, l *slog.Logger) error {
	storedToken, err := ts.Get(ctx, username)
	if err != nil {
		l.Error("log_id", logId, fmt.Sprintf("Failed to fetch stored token for '%v': %v", username, err))
		return HandleAuthError("AUTH003", logId, l)
	}

	if storedToken != providedToken {
		l.Debug("log_id", logId, "Provided token does not match stored token.")
		return HandleAuthError("AUTH003", logId, l)
	}
	return nil
}
