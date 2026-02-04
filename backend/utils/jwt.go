package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair holds both tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTManager struct {
	secret            []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	issuer             string
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWT manager with config
func NewJWTManager(secret string, accessExpiryHours, refreshExpiryDays int) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("JWT secret cannot be empty")
	}
	
	return &JWTManager{
		secret:             []byte(secret),
		accessTokenExpiry:  time.Duration(accessExpiryHours) * time.Hour,
		refreshTokenExpiry: time.Duration(refreshExpiryDays) * 24 * time.Hour,
		issuer:             "stackflow",
	}, nil
}

func (j *JWTManager) GenerateTokenPair(userID int, email, role string) (*TokenPair, error) {
	accessToken, err := j.generateToken(userID, email, role, "access", j.accessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := j.generateToken(userID, email, role, "refresh", j.refreshTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWTManager) generateToken(userID int, email, role, tokenType string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.validateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	if claims.Type != "access" {
		return nil, errors.New("invalid token type: expected access token")
	}
	
	return claims, nil
}

func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.validateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type: expected refresh token")
	}
	
	return claims, nil
}

func (j *JWTManager) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}