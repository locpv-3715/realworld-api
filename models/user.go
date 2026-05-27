package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Bio      string `gorm:"type:text"`
	Image    string `gorm:"type:varchar(255)"`
}

type Follow struct {
	FollowerID uint `gorm:"primaryKey;autoIncrement:false"`
	FolloweeID uint `gorm:"primaryKey;autoIncrement:false"`
}
