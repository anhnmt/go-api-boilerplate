package model

import (
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = (*User)(nil)

type User struct {
	BaseModel

	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"type:varchar(255);unique" json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `gorm:"type:varchar(255);not null;default:USER" json:"role"`
}

func (User) TableName() string {
	return "users"
}
