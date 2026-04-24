package handlers

import (
	"go-proj/config"
	"go-proj/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBoard(c *gin.Context) {
	var board models.Board
	c.BindJSON(&board)

	config.DB.Create(&board)
	c.JSON(http.StatusOK, board)
}

func GetBoards(c *gin.Context) {
	var boards []models.Board
	config.DB.Find(&boards)

	c.JSON(http.StatusOK, boards)
}
