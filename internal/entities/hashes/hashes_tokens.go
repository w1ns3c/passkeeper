package hashes

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/config"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// GenerateToken func for generate JWT auth token
func GenerateToken(userid string, secret string, lifetime time.Duration) (token string, err error) {
	tokenLife := time.Now().Add(lifetime)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenLife),
		},
		UserID: userid,
	})
	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// CheckToken func for validate JWT auth token
func CheckToken(tokenStr, secret string) (userID string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	// return user ID in readable format
	return claims.UserID, nil
}

// CheckToken func for validate JWT auth token
func ExtractUserID(token string) (userID string, err error) {
	claims := &Claims{}
	parser := jwt.NewParser()
	_, _, err = parser.ParseUnverified(token, claims)
	if err != nil {
		return "", err
	}

	// return user ID in readable format
	return claims.UserID, nil
}

func ExtractUserInfo(ctx context.Context) (token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", entities.ErrNoToken
	}

	tokens := md.Get(config.TokenHeader)
	if len(tokens) == 0 {
		return "", entities.ErrEmptyToken
	}
	token = tokens[0]

	return token, nil
}
