package jwtutils

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(refreshClaims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	return token.SignedString(secret)
}
