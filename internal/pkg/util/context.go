package util

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ctxAuthClaimsKey struct{}

func ExtractCtxClaims(ctx context.Context) (jwt.MapClaims, error) {
	claims, ok := ctx.Value(ctxAuthClaimsKey{}).(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed get claims")
	}

	return claims, nil
}

func SetCtxClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, ctxAuthClaimsKey{}, claims)
}

func ProvideCtx(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}
