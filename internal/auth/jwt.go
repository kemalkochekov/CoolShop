package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	claimsTimeLimit        = 5 * time.Minute
	refreshClaimsTimeLimit = 30 * time.Minute
)

type JwtGen interface {
	GenerateJWT(userID string) (string, string, error)
	ValidateToken(signedToken string) (*SignedDetails, error)
}

type SignedDetails struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var _ JwtGen = (*TokenGenerator)(nil)

type TokenGenerator struct {
	secretKey string
}

func NewTokenGenerator(secretKey string) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
	}
}

func (t *TokenGenerator) GenerateJWT(userID string) (string, string, error) {
	claims := SignedDetails{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(claimsTimeLimit)),
		},
	}
	refreshClaims := SignedDetails{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshClaimsTimeLimit)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(t.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (t *TokenGenerator) ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Errorf("the token is invalid!!! %w", err)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("the token is expired %w", err)
	}

	return claims, nil
}
