package models

import (
	"gorm.io/gorm"
	"html/template"
	"time"
)

type Novel struct {
	ID               uint   `gorm:"uniqueIndex; primaryKey"`
	SiteID           uint   `gorm:"not null"`
	Name             string `gorm:"index; size:255;not null"`
	template.URL     `gorm:"index; size:255;not null"`
	NumberOfChapters int `gorm:"index; not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
