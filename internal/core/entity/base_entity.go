package entity

import (
	"time"
)

type BaseEntity struct {
	ID        string    `gorm:"type:varchar(50);primaryKey;not null" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamp(6) with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp(6) with time zone;index" json:"updated_at"`
}
