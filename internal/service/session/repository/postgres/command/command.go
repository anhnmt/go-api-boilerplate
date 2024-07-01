package sessioncommand

import (
	"context"
	"time"

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

func (c *Command) UpdateIsRevoked(ctx context.Context, sessionId string, isRevoked bool, now time.Time) error {
	e := c.DB().Session

	_, err := c.db.WriteDB().Session.WithContext(ctx).Where(e.ID.Eq(sessionId)).
		Updates(map[string]interface{}{
			"is_revoked":   isRevoked,
			"last_seen_at": now,
			"updated_at":   now,
		})
	return err
}

func (c *Command) UpdateLastSeenAt(ctx context.Context, sessionId string, now time.Time) error {
	e := c.DB().Session

	_, err := c.db.WriteDB().Session.WithContext(ctx).Where(e.ID.Eq(sessionId)).
		Updates(map[string]interface{}{
			"last_seen_at": now,
			"updated_at":   now,
		})
	return err
}
