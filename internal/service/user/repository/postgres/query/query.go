package userquery

import (
	"context"

	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/model"
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

func (q *Query) ListUsers(ctx context.Context) ([]*model.User, error) {
	e := q.db.User

	return q.db.ReadDB().User.WithContext(ctx).
		Omit(e.Password).
		Find()
}

func (q *Query) GetByEmailWithPassword(ctx context.Context, email string) (*model.User, error) {
	e := q.db.User

	return q.db.ReadDB().User.WithContext(ctx).Select(e.ID, e.Name, e.Email, e.Role, e.Password).
		Where(e.Email.Eq(email)).
		First()
}

func (q *Query) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	e := q.db.User

	return q.db.ReadDB().User.WithContext(ctx).Omit(e.Password).
		Where(e.Email.Eq(email)).
		First()
}

func (q *Query) GetByID(ctx context.Context, ID string) (*model.User, error) {
	e := q.db.User

	return q.db.ReadDB().User.WithContext(ctx).
		Omit(e.Password).
		Where(e.ID.Eq(ID)).
		First()
}
