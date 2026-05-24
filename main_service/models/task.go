package models

import "time"

type Task struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Title     string     `json:"title" binding:"required"`
	Content   string     `json:"content"`
	Status    string     `json:"status" gorm:"default:'todo'"`     // todo, in-progress, done
	Priority  string     `json:"priority" gorm:"default:'medium'"` // low, medium, high
	BoardID   uint       `json:"board_id" binding:"required"`
	UserID    uint       `json:"user_id"`
	Deadline  *time.Time `json:"deadline"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
