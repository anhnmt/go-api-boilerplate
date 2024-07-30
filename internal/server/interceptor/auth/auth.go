package authinterceptor

import (
	"context"
	"slices"

	"github.com/casbin/casbin/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/spf13/cast"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	casbin       *casbin.Enforcer
	authBusiness *authbusiness.Business
}

type Param struct {
	fx.In

	Casbin       *casbin.Enforcer
	AuthBusiness *authbusiness.Business
}

func New(p Param) AuthInterceptor {
	return &authInterceptor{
		casbin:       p.Casbin,
		authBusiness: p.AuthBusiness,
	}
}

func (a *authInterceptor) AuthFunc() auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		fullMethod, ok := grpc.Method(ctx)
		if !(ok && slices.Contains(defaultGuardLists, fullMethod)) {
			return ctx, nil
		}

		rawToken, err := auth.AuthFromMD(ctx, util.TokenType)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed get token")
		}

		claims, err := a.authBusiness.ExtractClaims(rawToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		if claims[util.Typ] != util.TokenType {
			return nil, status.Errorf(codes.InvalidArgument, "invalid token")
		}

		sessionID := cast.ToString(claims[util.Sid])
		tokenID := cast.ToString(claims[util.Jti])
		err = a.authBusiness.CheckBlacklist(ctx, sessionID, tokenID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		enforce, err := a.casbin.Enforce(fullMethod, claims[util.Role])
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		if !enforce {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		ctx = util.SetCtxClaims(ctx, claims)
		return ctx, nil
	}
}
