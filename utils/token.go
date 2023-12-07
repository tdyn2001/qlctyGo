package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(ttl time.Duration, payload interface{}, secretKey string) (string, error) {

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, key string) (interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t_.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}
