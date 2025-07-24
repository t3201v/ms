package api

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	ErrMissingToken = errors.New("missing or invalid Authorization header")

	skipMethods = map[string]bool{
		"/grpc.health.v1.Health/Check":                                   true,
		"/grpc.health.v1.Health/List":                                    true,
		"/grpc.health.v1.Health/Watch":                                   true,
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
		"/resource.ResourceService/SayHello":                             true,
		"/resource/say":                                                  true,
		"/resource/say2":                                                 false,
		"/resource.ResourceService/SayHello2":                            false,
	}
)

// Load your public key (once, globally)
func LoadPublicKey() (*rsa.PublicKey, error) {
	data, err := os.ReadFile("public.key")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(data)
}

// Interceptor
func AuthInterceptor(publicKey *rsa.PublicKey) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if skipMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrMissingToken
		}

		// Extract Authorization: Bearer <token>
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 || !strings.HasPrefix(authHeaders[0], "Bearer ") {
			return nil, ErrMissingToken
		}

		rawToken := strings.TrimPrefix(authHeaders[0], "Bearer ")

		// Parse and verify the JWT
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		return handler(ctx, req)
	}
}
