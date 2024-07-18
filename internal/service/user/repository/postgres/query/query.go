package userquery

import (
	"context"

	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/model"
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
