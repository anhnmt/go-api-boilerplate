package base

import (
	"context"
)

func ProvideCtx(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}
