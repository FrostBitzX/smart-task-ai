package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const ProjectIDPrefix = "proj"

type Project struct {
	ID        uuid.UUID       `gorm:"type:char(36);primaryKey"`
	NodeID    *uuid.UUID      `gorm:"type:char(36)"`
	AccountID uuid.UUID       `gorm:"type:char(36);unique;not null"`
	Role      string          `gorm:"type:enum('owner','member');not null"`
	Name      string          `gorm:"type:varchar(255);not null"`
	Config    json.RawMessage `gorm:"type:jsonb"`
	CreatedAt time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time      `gorm:"default:null"`
}

func (Project) TableName() string {
	return "projects"
}
