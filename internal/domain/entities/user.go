package entities

import (
	"time"

	"gorm.io/gorm/schema"
)

var _ schema.Tabler = (*User)(nil)

type User struct {
	ID        string    `gorm:"type:varchar(50);primary_key;not null" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Email     string    `gorm:"type:varchar(255)" json:"email"`
	CreatedAt time.Time `gorm:"type:timestamp(6) with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp(6) with time zone;index" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
