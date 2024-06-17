package command

import (
	"gorm.io/gorm"
)

type Command struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Command {
	return &Command{
		db: db,
	}
}
