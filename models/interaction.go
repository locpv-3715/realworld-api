package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Body      string  `gorm:"type:text;not null"`
	ArticleID uint    `gorm:"not null"`
	Article   Article `gorm:"foreignKey:ArticleID"`
	AuthorID  uint    `gorm:"not null"`
	Author    User    `gorm:"foreignKey:AuthorID"`
}

type Favorite struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false"`
	ArticleID uint `gorm:"primaryKey;autoIncrement:false"`
}
