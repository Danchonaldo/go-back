package models

type Task struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"` // todo, in-progress, done
	BoardID uint   `json:"board_id"`
}
