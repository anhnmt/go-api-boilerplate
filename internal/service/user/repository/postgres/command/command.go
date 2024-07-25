package usercommand

import (
	"context"

	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/model"
)

type Command struct {
	db *gormgen.Query
}

type Params struct {
	fx.In

	DB *gormgen.Query
}

func New(p Params) *Command {
	return &Command{
		db: p.DB,
	}
}

func (c *Command) Create(ctx context.Context, user *model.User) error {
	return c.db.WriteDB().User.WithContext(ctx).Create(user)
}
