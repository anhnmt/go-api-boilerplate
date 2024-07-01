package credentialentity

import (
	"time"

	"gorm.io/gorm/schema"

	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/core/entity"
)

var _ schema.Tabler = (*Credential)(nil)

type Credential struct {
	entity.BaseEntity

	UserID         string          `gorm:"type:varchar(50);not null;index" json:"user_id"`
	Type           string          `gorm:"type:varchar(255)" json:"type"`
	SecretData     *SecretData     `gorm:"type:jsonb;serializer:json;default:NULL" json:"secret_data"`
	CredentialData *CredentialData `gorm:"type:jsonb;serializer:json;default:NULL" json:"credential_data"`
	ExpiredAt      time.Time       `gorm:"type:timestamp(6) with time zone" json:"expired_at"`
}

func (Credential) TableName() string {
	return "credentials"
}

type SecretData struct {
}

type CredentialData struct {
}
