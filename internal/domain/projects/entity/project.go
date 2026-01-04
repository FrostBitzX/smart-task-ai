package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const ProjectIDPrefix = "proj"

type Project struct {
	ID        uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey"`
	NodeID    *uuid.UUID      `json:"nodeId" gorm:"type:char(36)"`
	AccountID uuid.UUID       `json:"accountId" gorm:"type:char(36);not null"`
	Role      string          `json:"role" gorm:"type:enum('owner','member');not null"`
	Name      string          `json:"name" gorm:"type:varchar(255);not null"`
	Config    json.RawMessage `json:"config" gorm:"type:jsonb"`
	CreatedAt time.Time       `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time      `json:"deletedAt" gorm:"default:null"`
}

func (Project) TableName() string {
	return "projects"
}
