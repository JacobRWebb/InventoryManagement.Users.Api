package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CommonBase represents a set of common fields used across all models in the system.
// It provides standardized attributes for entity identification and tracking creation/modification times.
type CommonBase struct {
	// Id is a unique identifier for the entity.
	// It is automatically generated as a UUID v4 upon creation.
	Id uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	// CreatedAt records the UTC timestamp when the entity was created.
	// It is automatically set to the current time upon entity creation.
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	// UpdatedAt records the UTC timestamp when the entity was last modified.
	// It is automatically updated to the current time whenever the entity is modified.
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (base *CommonBase) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()

	if base.CreatedAt.IsZero() {
		base.CreatedAt = now
	}

	if base.UpdatedAt.IsZero() {
		base.UpdatedAt = now
	}

	return nil
}

func (base *CommonBase) BeforeUpdate(tx *gorm.DB) error {
	base.UpdatedAt = time.Now().UTC()

	return nil
}
