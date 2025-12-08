package entity

import "time"

type Account struct {
	ID        string     `gorm:"type:char(36);primaryKey"`
	NodeID    string     `gorm:"type:char(36);not null"`
	Username  string     `gorm:"type:varchar(100);unique;not null"`
	Email     string     `gorm:"type:varchar(255);unique;not null"`
	Password  string     `gorm:"type:varchar(255);not null"`
	State     string     `gorm:"type:enum('pending','active','inactive');not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"` // soft delete
}

func (Account) TableName() string {
	return "accounts"
}
