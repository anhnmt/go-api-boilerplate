package entities

import (
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = (*Device)(nil)

type Device struct {
	BaseEntity

	UserID      string `gorm:"type:varchar(50);not null;index" json:"user_id"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Fingerprint string `gorm:"type:varchar(255)" json:"fingerprint"`
	UserAgent   string `gorm:"type:varchar(255)" json:"user_agent"`
	IpAddress   string `gorm:"type:varchar(255)" json:"user_agent"`
	Secret      string `gorm:"type:varchar(255)" json:"secret"`
}

func (Device) TableName() string {
	return "devices"
}
