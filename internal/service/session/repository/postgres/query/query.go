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

func (q *Query) FindByUserIDAndSessionID(ctx context.Context, userID, sessionID string) ([]*pb.ActiveSessions, error) {
	return q.db.ReadDB().Session.WithContext(ctx).FindByUserIDAndSessionID(userID, sessionID)
}

func (q *Query) FindByUserIDWithoutSessionID(ctx context.Context, userID, sessionID string) ([]string, error) {
	return q.db.ReadDB().Session.WithContext(ctx).FindByUserIDWithoutSessionID(userID, sessionID)
}
