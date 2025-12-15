package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	NodeID    *uuid.UUID `gorm:"type:char(36)"`
	Username  string     `gorm:"type:varchar(100);unique;not null"`
	Email     string     `gorm:"type:varchar(255);unique;not null"`
	Password  string     `gorm:"type:varchar(255);not null"`
	State     string     `gorm:"type:enum('active','inactive');not null;default:'active'"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"` // soft delete
}

func (Account) TableName() string {
	return "accounts"
}
