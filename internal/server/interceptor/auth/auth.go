package authinterceptor

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
)

var defaultGuardLists = []string{
	"/auth.v1.AuthService/Info",
	"/auth.v1.AuthService/RevokeToken",
}

type AuthInterceptor interface {
	AuthFunc() auth.AuthFunc
}

type authInterceptor struct {
	authBusiness *authbusiness.Business
}

func New(
	authBusiness *authbusiness.Business,
) AuthInterceptor {
	return &authInterceptor{
		authBusiness: authBusiness,
	}
}

func (a *authInterceptor) AuthFunc() auth.AuthFunc {
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

			sessionId := cast.ToString(claims[jwtutils.Sid])
			tokenId := cast.ToString(claims[jwtutils.Jti])
			err = a.authBusiness.CheckBlacklist(ctx, sessionId, tokenId)
			if err != nil {
				return nil, err
			}

			log.Info().Any("claims", claims).Msg("get claims")
		}

		return ctx, nil
	}
}

func (a *authInterceptor) checkFullMethod(ctx context.Context) bool {
	fullMethod, ok := grpc.Method(ctx)
	if !ok {
		return false
	}

	return lo.Contains(defaultGuardLists, fullMethod)
}
