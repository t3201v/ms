package api

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func issueJWT(sub string, scopes []string) (string, error) {
	privateKeyData, err := os.ReadFile("../private.key")
	if err != nil {
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   sub,
		"iss":   "https://ms-issuer.com",
		"scope": strings.Join(scopes, " "),
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(privateKey)
}
