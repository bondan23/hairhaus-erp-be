package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// AuditLog records significant system actions for audit trail.
type AuditLog struct {
	BaseModel
	Action     string    `gorm:"not null" json:"action"`
	EntityType string    `gorm:"not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`

	PerformedBy uuid.UUID `gorm:"type:uuid;not null" json:"performed_by"`

	Metadata json.RawMessage `gorm:"type:jsonb" json:"metadata,omitempty"`
}
