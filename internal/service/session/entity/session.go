package sessionentity

import (
	"time"

	"gorm.io/gorm/schema"

	"github.com/anhnmt/go-api-boilerplate/internal/core/entity"
)

var _ schema.Tabler = (*Session)(nil)

type Session struct {
	entity.BaseEntity

	UserID      string     `gorm:"type:varchar(50);not null;index" json:"user_id"`
	Fingerprint string     `gorm:"type:varchar(50)" json:"fingerprint"`
	UserAgent   string     `gorm:"type:varchar(255)" json:"user_agent"`
	DeviceType  string     `gorm:"type:varchar(50)" json:"device_type"`
	OS          string     `gorm:"type:varchar(50)" json:"os"`
	Browser     string     `gorm:"type:varchar(50)" json:"browser"`
	Device      string     `gorm:"type:varchar(50)" json:"device"`
	IpAddress   string     `gorm:"type:varchar(100)" json:"ip_address"`
	IsRevoked   bool       `json:"is_revoked"`
	LastSeenAt  *time.Time `gorm:"type:timestamp(6) with time zone" json:"last_seen_at"`
	// ExpiredAt   *time.Time `gorm:"type:timestamp(6) with time zone" json:"expired_at"`
}

func (Session) TableName() string {
	return "sessions"
}
