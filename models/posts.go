package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	User_id  uint   `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	IsPublic bool   `gorm:"default:true" json:"isPublic"`
}
