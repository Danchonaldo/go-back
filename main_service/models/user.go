package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" gorm:"unique" binding:"required,email"`
	Password  string    `json:"-"`
	Role      string    `json:"role" gorm:"default:'user'"` // user, moderator, admin
	CreatedAt time.Time `json:"created_at"`
}
