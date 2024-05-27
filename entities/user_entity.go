package entities

import (
	"gorm.io/gorm"
	"time"
)

type UserEntity struct {
	IdUser     uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	IdTelegram int64          `gorm:"not null"`
	Admin      bool           `gorm:"default:false"`
	Name       string         `gorm:"not null"`
}

func (UserEntity) TableName() string {
	return "users"
}
