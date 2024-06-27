package sessionquery

import (
	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
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
