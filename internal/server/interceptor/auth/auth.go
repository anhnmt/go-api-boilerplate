package authinterceptor

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/util"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
)

var defaultGuardLists = []string{
	"/auth.v1.AuthService/Info",
	"/auth.v1.AuthService/ActiveSessions",
	"/user.v1.UserService/ListUsers",
}

type AuthInterceptor interface {
	AuthFunc() auth.AuthFunc
}

type authInterceptor struct {
	authBusiness *authbusiness.Business
}

type Param struct {
	fx.In

	AuthBusiness *authbusiness.Business
}

func New(p Param) AuthInterceptor {
	return &authInterceptor{
		authBusiness: p.AuthBusiness,
	}
}

func (a *authInterceptor) AuthFunc() auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if !a.checkFullMethod(ctx) {
			return ctx, nil
		}

		rawToken, err := auth.AuthFromMD(ctx, util.TokenType)
		if err != nil {
			return nil, fmt.Errorf("failed get token")
		}

		claims, err := a.authBusiness.ExtractClaims(rawToken)
		if err != nil {
			return nil, err
		}

		ctx = util.SetCtxClaims(ctx, claims)
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
