package model

import "github.com/jinzhu/gorm"

// 用户的模型
type User struct {
	gorm.Model
	Name  string `gorm:"type:varchar(20);not null"`
	Phone string `gorm:"varchar(11);not null;unique"`
	Pwd   string `gorm:"size:255;not null"`
}
