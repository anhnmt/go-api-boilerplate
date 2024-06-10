package entities

import (
	"time"

	"gorm.io/gorm/schema"
)

var _ schema.Tabler = (*Session)(nil)

type Session struct {
	BaseEntity

	DeviceID  string    `gorm:"type:varchar(50);not null;index" json:"device_id"`
	Token     string    `gorm:"type:varchar(255)" json:"token"`
	IsRevoked bool      `json:"is_revoked"`
	ExpiredAt time.Time `gorm:"type:timestamp(6) with time zone" json:"expired_at"`
}

func (Session) TableName() string {
	return "sessions"
}