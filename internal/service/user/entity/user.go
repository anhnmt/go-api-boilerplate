package userentity

import (
	"gorm.io/gorm/schema"

	"github.com/anhnmt/go-api-boilerplate/internal/core/entities"
)

var _ schema.Tabler = (*User)(nil)

type User struct {
	entities.BaseEntity

	Name  string `gorm:"type:varchar(255)" json:"name"`
	Email string `gorm:"type:varchar(255)" json:"email"`
}

func (User) TableName() string {
	return "users"
}
