// Package jwt provides JWT (JSON Web Token) functionality for authentication and authorization.
// It includes token creation, verification, and claims extraction utilities.
package jwt

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// contextKey represents a type for context keys to avoid collisions.
type contextKey string

// ContextClaims is the context key used to store JWT claims in HTTP request context.
const (
	ContextClaims contextKey = "claims"
)

// JWTClaims represents the claims structure for JWT tokens.
type JWTClaims struct {
	jwtlib.RegisteredClaims
}

// VerifyToken validates a JWT token string using the provided secret key.
func VerifyToken(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return &JWTClaims{}, err
	}

	if !token.Valid {
		return &JWTClaims{}, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return &JWTClaims{}, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

// GetJwtClaims extracts JWT claims from the HTTP request context.
func GetJwtClaims(r *http.Request) (*JWTClaims, error) {
	claims, ok := r.Context().Value(ContextClaims).(*JWTClaims)
	if !ok {
		return &JWTClaims{}, fmt.Errorf("failed to get validated claims")
	}
	return claims, nil
}

// CreateJwtToken creates a new JWT token with the specified user ID and expiration time.
func CreateJwtToken(JwtSecretToken string, userID int) (string, error) {
	claims := JWTClaims{
		jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Duration(3600 * time.Second))),
			Subject:   strconv.Itoa(userID),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(JwtSecretToken))
	if err != nil {
		return "", err
	}
	return t, nil
}
