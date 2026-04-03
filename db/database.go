package db

import (
	"go-proj/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=Dancho2005 dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to DB")
	}

	DB = db

	DB.AutoMigrate(&models.User{}, &models.Board{}, &models.Task{})
}
