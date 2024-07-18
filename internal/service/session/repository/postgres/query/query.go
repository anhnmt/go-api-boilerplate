package sessionquery

import (
	"context"

	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	"github.com/anhnmt/go-api-boilerplate/gen/pb"
)

type Query struct {
	db *gormgen.Query
}

type Params struct {
	fx.In

	DB *gormgen.Query
}

func New(p Params) *Query {
	return &Query{
		db: p.DB,
	}
}

func (q *Query) FindByUserIdAndSessionId(ctx context.Context, userId, sessionId string) ([]*pb.ActiveSessions, error) {
	return q.db.ReadDB().Session.WithContext(ctx).FindByUserIdAndSessionId(userId, sessionId)
}

func (q *Query) FindByUserIdWithoutSessionId(ctx context.Context, userId, sessionId string) ([]string, error) {
	return q.db.ReadDB().Session.WithContext(ctx).FindByUserIdWithoutSessionId(userId, sessionId)
}
