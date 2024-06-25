package userquery

import (
	"context"

	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
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

func (q *Query) ListUsers(ctx context.Context) ([]*userentity.User, error) {
	e := q.DB().User

	return q.db.ReadDB().User.WithContext(ctx).Select(e.ID, e.Name, e.Email, e.CreatedAt, e.UpdatedAt).Find()
}
