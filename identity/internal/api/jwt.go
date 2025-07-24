package api

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func issueJWT() (string, error) {
	privateKeyData, err := os.ReadFile("../private.key")
	if err != nil {
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "user123",
		"iss": "https://my-issuer.example.com",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(privateKey)
}
