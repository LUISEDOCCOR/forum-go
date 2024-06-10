package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	UserId   float64 `json:"user_id"`
	Author   string  `json:"author"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	IsPublic bool    `gorm:"default:true" json:"isPublic"`
}
