package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

var _ callbacks.BeforeCreateInterface = (*BaseModel)(nil)

// BaseModel contains common columns for all tables.
type BaseModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;not null" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamp(6) with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp(6) with time zone;index" json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *BaseModel) BeforeCreate(*gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}

	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}

	return nil
}
