package models

import "time"

type Board struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}
