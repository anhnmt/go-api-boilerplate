package sessioncommand

import (
	"context"

	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	sessionentity "github.com/anhnmt/go-api-boilerplate/internal/service/session/entity"
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

func (c *Command) Create(ctx context.Context, session *sessionentity.Session) error {
	return c.db.WriteDB().Session.WithContext(ctx).Create(session)
}
