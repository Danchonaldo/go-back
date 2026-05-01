package handlers

import (
	"go-proj/main_service/db"
	"go-proj/main_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBoard(c *gin.Context) {
	var board models.Board
	c.BindJSON(&board)

	db.DB.Create(&board)
	c.JSON(http.StatusOK, board)
}

func GetBoards(c *gin.Context) {
	var boards []models.Board
	db.DB.Find(&boards)

	c.JSON(http.StatusOK, boards)
}
