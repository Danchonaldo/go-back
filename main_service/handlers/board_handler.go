package handlers

import (
	"main_service/db"
	"main_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBoard(c *gin.Context) {
	var board models.Board

	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	board.UserID = userID.(uint)

	if err := db.DB.Create(&board).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board"})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func GetBoards(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var boards []models.Board

	if role == "admin" || role == "moderator" {
		db.DB.Find(&boards)
	} else {
		db.DB.Where("user_id = ?", userID).Find(&boards)
	}

	c.JSON(http.StatusOK, boards)
}

func GetBoardByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var board models.Board

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&board).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}

	c.JSON(http.StatusOK, board)
}

func UpdateBoard(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var board models.Board

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&board).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}

	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Save(&board)
	c.JSON(http.StatusOK, board)
}

func DeleteBoard(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var board models.Board

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&board).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}

	db.DB.Delete(&board)
	c.JSON(http.StatusOK, gin.H{"message": "Board deleted successfully"})
}
