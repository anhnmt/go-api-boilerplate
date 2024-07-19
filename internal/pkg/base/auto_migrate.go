package base

import (
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/internal/model"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Session{},
	)
}
