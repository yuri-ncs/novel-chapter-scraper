package models

import "time"

//generate a user struct that have their respectives chat ids and their novels list ids to be notified

type User struct {
	ID        uint    `gorm:"uniqueIndex; primaryKey"`
	ChatID    int64   `gorm:"uniqueIndex; not null"`
	Novels    []Novel `gorm:"many2many:user_novels;"`
	CreatedAt time.Time
	DeletedAt time.Time
}
