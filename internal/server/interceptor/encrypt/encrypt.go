package encryptinterceptor

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/grpc"
)

var defaultGuardLists = []string{
	"/auth.v1.AuthService/Encrypt",
}

type EncryptInterceptor interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
}

type encryptInterceptor struct {
}

func New() EncryptInterceptor {
	return &encryptInterceptor{}
}

func (e *encryptInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !e.checkFullMethod(ctx, info) {
			return handler(ctx, req)
		}

		return handler(ctx, req)
	}
}

func (e *encryptInterceptor) checkFullMethod(ctx context.Context, info *grpc.UnaryServerInfo) bool {
	fullMethod, ok := grpc.Method(ctx)
	if !ok {
		return lo.Contains(defaultGuardLists, info.FullMethod)
	}

	return lo.Contains(defaultGuardLists, fullMethod)
}
