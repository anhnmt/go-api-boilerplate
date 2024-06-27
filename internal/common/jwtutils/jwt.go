package jwtutils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(refreshClaims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if val, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || val.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return keyFunc(token)
	})
}
