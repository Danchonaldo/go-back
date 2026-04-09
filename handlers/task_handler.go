package handlers

import (
	"go-proj/config"
	"go-proj/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task
	c.BindJSON(&task)

	config.DB.Create(&task)
	c.JSON(http.StatusOK, task)
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	config.DB.Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	config.DB.First(&task, id)
	c.BindJSON(&task)
	config.DB.Save(&task)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	config.DB.Delete(&models.Task{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func GetTasksByBoard(c *gin.Context) {
	boardID := c.Param("id")
	var tasks []models.Task

	config.DB.Where("board_id = ?", boardID).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}
