package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"proyecto.com/m/models"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Group  string `json:"group"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) CreateJWT(user models.User) (TokenPair, error) {
	return CreateJWT(user)
}

func (s *AuthService) ValidateJWT(tokenString string) (JWTClaims, error) {
	return ValidateJWT(tokenString)
}

func CreateJWT(user models.User) (TokenPair, error) {
	accessSecret := os.Getenv("JWT_SECRET")
	if accessSecret == "" {
		accessSecret = "dev-secret-change-me"
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = accessSecret
	}

	now := time.Now()
	accessClaims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		Group:  user.Group,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	refreshClaims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		Group:  user.Group,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(accessSecret))
	if err != nil {
		return TokenPair{}, fmt.Errorf("could not sign access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(refreshSecret))
	if err != nil {
		return TokenPair{}, fmt.Errorf("could not sign refresh token: %w", err)
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func ValidateJWT(tokenString string) (JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	claims := JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return JWTClaims{}, ErrInvalidToken
	}

	if !token.Valid {
		return JWTClaims{}, ErrInvalidToken
	}

	if claims.UserID == "" || claims.Role == "" || claims.Group == "" {
		return JWTClaims{}, ErrInvalidToken
	}

	return claims, nil
}
