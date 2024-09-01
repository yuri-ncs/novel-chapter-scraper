package models

import (
	"gorm.io/gorm"
	"time"
)

type Site struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:255;not null"`
	DefaultURL string `gorm:"size:255;not null;unique"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
