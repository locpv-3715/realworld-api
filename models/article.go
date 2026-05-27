package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Slug        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Title       string `gorm:"type:varchar(255);not null"`
	Description string
	Body        string
	AuthorID    uint
	Author      User  `gorm:"foreignKey:AuthorID"`
	Tags        []Tag `gorm:"many2many:article_tags;"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);uniqueIndex;not null"`
}
