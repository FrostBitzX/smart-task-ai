package entity

import (
	"time"

	"github.com/google/uuid"
)

const ProfileIDPrefix = "prof"

type Profile struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	NodeID     uuid.UUID `gorm:"type:char(36);not null;index"`
	AccountID  uuid.UUID `gorm:"type:char(36);unique;not null"`
	FirstName  string    `gorm:"type:varchar(100);not null"`
	LastName   string    `gorm:"type:varchar(100);not null"`
	Nickname   *string   `gorm:"type:varchar(50)"`
	AvatarPath *string   `gorm:"type:varchar(500)"`
	State      string    `gorm:"type:enum('active','inactive');not null"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Profile) TableName() string {
	return "profiles"
}
