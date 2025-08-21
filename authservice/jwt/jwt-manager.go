package jwt

import (
	"authservice/models"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var logger *zap.Logger

type JWTManager struct {
	secretKey            string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	UserEmail string `json:"email"`
	Type      string `json:"type"`
}

func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		tokenDuration:        accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (jwtManager *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtManager.tokenDuration).Unix(),
		},
		UserEmail: user.Email,
		Type:      "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtManager.secretKey))
}

func (jwtManager *JWTManager) GenerateRefreshToken(user *models.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtManager.refreshTokenDuration).Unix(),
		},
		UserEmail: user.Email,
		Type:      "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtManager.secretKey))
}

func (jwtManager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(jwtManager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
