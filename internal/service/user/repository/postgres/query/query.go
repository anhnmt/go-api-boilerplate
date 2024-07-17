package userquery

import (
	"context"

	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
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

func (q *Query) GetByEmailWithPassword(ctx context.Context, email string) (*userentity.User, error) {
	e := q.DB().User

	return q.db.ReadDB().User.WithContext(ctx).Select(e.ID, e.Name, e.Email, e.Password).
		Where(e.Email.Eq(email)).
		First()
}

func (q *Query) GetByEmail(ctx context.Context, email string) (*userentity.User, error) {
	e := q.DB().User

	return q.db.ReadDB().User.WithContext(ctx).Omit(e.Password).
		Where(e.Email.Eq(email)).
		First()
}

func (q *Query) GetByID(ctx context.Context, id string) (*userentity.User, error) {
	e := q.DB().User

	return q.db.ReadDB().User.WithContext(ctx).Omit(e.Password).
		Where(e.ID.Eq(id)).
		First()
}
