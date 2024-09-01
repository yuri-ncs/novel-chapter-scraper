package models

import (
	"gorm.io/gorm"
	"time"
)

type Chapter struct {
	ID        uint   `gorm:"primaryKey"`
	NovelID   uint   `gorm:"not null"`
	Number    int    `gorm:"not null"`
	Title     string `gorm:"size:255;not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
