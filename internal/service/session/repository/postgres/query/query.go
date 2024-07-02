package sessionquery

import (
	"context"

	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Query struct {
	db *gormgen.Query
}

func New(db *gormgen.Query) *Query {
	return &Query{
		db: db,
	}
}

func (q *Query) DB() *gormgen.Query {
	return q.db
}

func (q *Query) FindByUserIdAndSessionId(ctx context.Context, userId, sessionId string) ([]*pb.ActiveSessions, error) {
	return q.db.ReadDB().Session.WithContext(ctx).FindByUserIdAndSessionId(userId, sessionId)
}

func (q *Query) FindByUserIdWithoutSessionId(ctx context.Context, userId, sessionId string) ([]string, error) {
	return q.db.ReadDB().Session.WithContext(ctx).FindByUserIdWithoutSessionId(userId, sessionId)
}
