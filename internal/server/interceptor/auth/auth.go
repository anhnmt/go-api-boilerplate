package authinterceptor

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
)

var defaultGuardLists = []string{
	"/auth.v1.AuthService/Info",
	"/auth.v1.AuthService/RevokeToken",
}

type AuthInterceptor struct {
	authBusiness *authbusiness.Business
}

func New(
	authBusiness *authbusiness.Business,
) *AuthInterceptor {
	return &AuthInterceptor{
		authBusiness: authBusiness,
	}
}

func (a *AuthInterceptor) AuthFunc() auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if a.checkFullMethod(ctx) {
			rawToken, err := auth.AuthFromMD(ctx, jwtutils.TokenType)
			if err != nil {
				return nil, fmt.Errorf("failed get token")
			}

			claims, err := a.authBusiness.ExtractClaims(rawToken)
			if err != nil {
				return nil, err
			}

			log.Info().Any("claims", claims).Msg("get claims")
		}

		return ctx, nil
	}
}

func (a *AuthInterceptor) checkFullMethod(ctx context.Context) bool {
	fullMethod, ok := grpc.Method(ctx)
	if !ok {
		return false
	}

	return lo.Contains(defaultGuardLists, fullMethod)
}
