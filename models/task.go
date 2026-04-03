package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
	BoardID uint   `json:"board_id"`
}
