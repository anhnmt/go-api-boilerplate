package postgres

import (
	"database/sql"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func GormDB(p *Postgres) *gorm.DB {
	return p.DB
}

func SqlDB(p *Postgres) (*sql.DB, error) {
	return p.SqlDB()
}

// Module provided to fx
var Module = fx.Module("postgres", fx.Provide(
	New,
	GormDB,
	SqlDB,
))
