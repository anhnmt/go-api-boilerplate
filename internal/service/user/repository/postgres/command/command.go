package usercommand

import (
	"context"

	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/model"
)

type Command struct {
	db *gormgen.Query
}

func New(db *gormgen.Query) *Command {
	return &Command{
		db: db,
	}
}

func (c *Command) DB() *gormgen.Query {
	return c.db
}

func (c *Command) Create(ctx context.Context, user *userentity.User) error {
	return c.db.WriteDB().User.WithContext(ctx).Create(user)
}
