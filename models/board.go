package models

type Board struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	UserID uint   `json:"user_id"`
}
