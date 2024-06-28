package sessioncommand

import (
	"context"

	"gorm.io/gorm/clause"

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

func (c *Command) CreateOnConflict(ctx context.Context, session *sessionentity.Session) error {
	return c.db.WriteDB().Session.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_seen_at", "expires_at", "updated_at"}),
	}).Create(session)
}
