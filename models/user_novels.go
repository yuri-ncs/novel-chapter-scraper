package models

type UserNovel struct {
	UserID  uint `gorm:"primaryKey; foreignKey:User.ID"`
	NovelID uint `gorm:"primaryKey; foreignKey:Novel.ID"`
}
