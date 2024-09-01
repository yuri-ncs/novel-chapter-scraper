package models

import (
	"gorm.io/gorm"
	"html/template"
	"time"
)

type Novel struct {
	ID               uint   `gorm:"primaryKey"`
	SiteID           uint   `gorm:"not null"`
	Name             string `gorm:"size:255;not null"`
	template.URL     `gorm:"size:255;not null"`
	NumberOfChapters int `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
