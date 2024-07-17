package userentity

import (
	"gorm.io/gorm/schema"

	"github.com/anhnmt/go-api-boilerplate/internal/core/entity"
)

var _ schema.Tabler = (*User)(nil)

type User struct {
	entity.BaseEntity

	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"type:varchar(255);unique" json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
