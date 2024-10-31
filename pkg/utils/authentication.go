package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationDetails struct {
	Username string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

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

	if err != nil {
		// logger.Info(err.Error())
	}

	return token, err
}

func HashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
