package handlers

import (
	"main_service/db"
	"main_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	resp, err := SendNotification("task created")
	if err != nil {
		c.JSON(500, gin.H{"error": "notification failed"})
		return
	}

	c.JSON(200, gin.H{
		"message": "task created",
		"notify":  resp,
	})
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	db.DB.Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := db.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	db.DB.First(&task, id)
	c.BindJSON(&task)
	db.DB.Save(&task)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	db.DB.Delete(&models.Task{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func GetTasksByBoard(c *gin.Context) {
	boardID := c.Param("id")
	var tasks []models.Task

	db.DB.Where("board_id = ?", boardID).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}
