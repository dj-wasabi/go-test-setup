package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationDetails struct {
	Username string
	jwt.StandardClaims
}

var (
	ErrNoAuthHeader             = errors.New("authorization header is missing")
	ErrInvalidAuthHeader        = errors.New("authorization header is malformed")
	ErrInvalidToken             = errors.New("token is invalid")
	ErrTokenExpired             = errors.New("the token is expired.")
	SECRET_KEY           string = os.Getenv("SECRET_KEY")
)

func ValidatePassword(providedpassword, storedpassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(storedpassword), []byte(providedpassword))
	if err != nil {
		return false
	}
	return true
}

func GenerateToken(username string) (signedToken string, err error) {

	claims := &AuthenticationDetails{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	return token, err
}

func GetBearerToken(log *slog.Logger, req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	log.Info("bladiebla", "token info", fmt.Sprintf("%v", authHdr))
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}

	prefix := "Bearer: "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}

	return strings.TrimPrefix(authHdr, prefix), nil
}

// func checkTokenClaims(expectedClaims []string, t AuthenticationDetails) error {
// 	claims := strings.Split(t.Scope, " ")
// 	claimsMap := make(map[string]bool, len(claims))
// 	for _, c := range claims {
// 		claimsMap[c] = true
// 	}

// 	for _, e := range expectedClaims {
// 		if !claimsMap[e] {
// 			return ErrClaimsInvalid
// 		}
// 	}

// 	return nil
// }

func HashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ValidateToken(l *slog.Logger, signedToken string) (claims *AuthenticationDetails, msg error) {
	l.Info("Validating the provided token.")
	token, err := jwt.ParseWithClaims(
		signedToken,
		&AuthenticationDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthenticationDetails)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrTokenExpired
	}

	return claims, msg
}
