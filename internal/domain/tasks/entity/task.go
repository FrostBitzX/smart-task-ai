package entity

import (
	"time"

	"github.com/google/uuid"
)

const TaskIDPrefix = "tsk"

type Task struct {
	ID             uuid.UUID  `json:"id" gorm:"column:id"`
	NodeID         *uuid.UUID `json:"nodeId" gorm:"column:node_id"`
	ProjectID      uuid.UUID  `json:"projectId" gorm:"column:project_id"`
	Name           string     `json:"name" gorm:"column:name"`
	Description    *string    `json:"description" gorm:"column:description"`
	Priority       string     `json:"priority" gorm:"column:priority"`
	StartDateTime  *string    `json:"startDateTime" gorm:"column:start_datetime"`
	EndDateTime    *string    `json:"endDateTime" gorm:"column:end_datetime"`
	Location       *string    `json:"location" gorm:"column:location"`
	RecurringDays  *int       `json:"recurringDays" gorm:"column:recurring_days"`
	RecurringUntil *string    `json:"recurringUntil" gorm:"column:recurring_until"`
	Status         string     `json:"status" gorm:"column:status"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

func (Task) TableName() string {
	return "tasks"
}
