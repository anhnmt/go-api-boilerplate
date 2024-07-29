package base

import (
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/internal/model"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/permission"
)

func GormAutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Session{},
	)
}

func PermissionAutoMigrate(permission *permission.Permission) error {
	return permission.AutoMigrate()
}
