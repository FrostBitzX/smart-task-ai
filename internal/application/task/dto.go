package task

import (
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/common"
)

type CreateTaskRequest struct {
	Name           string  `json:"name" validate:"required"`
	Description    *string `json:"description"`
	Priority       string  `json:"priority" validate:"required"`
	StartDateTime  *string `json:"start_datetime"`
	EndDateTime    *string `json:"end_datetime"`
	Location       *string `json:"location"`
	RecurringDays  *int    `json:"recurring_days"`
	RecurringUntil *string `json:"recurring_until"`
}

type CreateTaskResponse struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	Priority       string  `json:"priority"`
	StartDateTime  *string `json:"start_datetime,omitempty"`
	EndDateTime    *string `json:"end_datetime,omitempty"`
	Location       *string `json:"location,omitempty"`
	RecurringDays  *int    `json:"recurring_days,omitempty"`
	RecurringUntil *string `json:"recurring_until,omitempty"`
}

type GetTaskByIDResponse struct {
	ID             string    `json:"id"`
	Status         string    `json:"status"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	Priority       string    `json:"priority"`
	StartDateTime  *string   `json:"start_datetime,omitempty"`
	EndDateTime    *string   `json:"end_datetime,omitempty"`
	Location       *string   `json:"location,omitempty"`
	RecurringDays  *int      `json:"recurring_days,omitempty"`
	RecurringUntil *string   `json:"recurring_until,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ListTasksByProjectResponse struct {
	Items      []GetTaskByIDResponse `json:"items"`
	Pagination common.Pagination     `json:"pagination"`
}

type UpdateTaskRequest struct {
	Name           string  `json:"name" validate:"required"`
	Description    *string `json:"description"`
	Priority       string  `json:"priority" validate:"required"`
	StartDateTime  *string `json:"start_datetime"`
	EndDateTime    *string `json:"end_datetime"`
	Location       *string `json:"location"`
	RecurringDays  *int    `json:"recurring_days"`
	RecurringUntil *string `json:"recurring_until"`
}

type UpdateTaskResponse struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	Priority       string  `json:"priority"`
	StartDateTime  *string `json:"start_datetime,omitempty"`
	EndDateTime    *string `json:"end_datetime,omitempty"`
	Location       *string `json:"location,omitempty"`
	RecurringDays  *int    `json:"recurring_days,omitempty"`
	RecurringUntil *string `json:"recurring_until,omitempty"`
}
