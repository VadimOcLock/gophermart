package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(subject string, expiresIn time.Time, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"exp": expiresIn.Unix(),
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
