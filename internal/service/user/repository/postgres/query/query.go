package query

import (
	"gorm.io/gorm"
)

type Query struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Query {
	return &Query{
		db: db,
	}
}
