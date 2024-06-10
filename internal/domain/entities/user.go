package entities

import (
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = (*User)(nil)

type User struct {
	BaseEntity

	Name  string `gorm:"type:varchar(255)" json:"name"`
	Email string `gorm:"type:varchar(255)" json:"email"`
}

func (User) TableName() string {
	return "users"
}
