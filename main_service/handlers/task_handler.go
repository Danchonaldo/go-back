package handlers

import (
	"main_service/db"
	"main_service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var input struct {
		Title    string  `json:"title" binding:"required"`
		Content  string  `json:"content"`
		Status   string  `json:"status"`
		Priority string  `json:"priority"`
		BoardID  uint    `json:"board_id" binding:"required"`
		Deadline *string `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	task := models.Task{
		Title:    input.Title,
		Content:  input.Content,
		Status:   input.Status,
		Priority: input.Priority,
		BoardID:  input.BoardID,
		UserID:   userID.(uint),
	}

	if task.Status == "" {
		task.Status = "todo"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}

	if input.Deadline != nil && *input.Deadline != "" {
		t, err := time.Parse(time.RFC3339, *input.Deadline)
		if err != nil {
			// try date-only format
			t, err = time.Parse("2006-01-02", *input.Deadline)
		}
		if err == nil {
			task.Deadline = &t
		}
	}

	if err := db.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	go func() {
		SendNotification("Task created: " + task.Title)
	}()

	c.JSON(http.StatusCreated, task)
}

func GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	var tasks []models.Task

	db.DB.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var task models.Task

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var task models.Task

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		Title    string  `json:"title"`
		Content  string  `json:"content"`
		Status   string  `json:"status"`
		Priority string  `json:"priority"`
		BoardID  uint    `json:"board_id"`
		Deadline *string `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	task.Content = input.Content
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Priority != "" {
		task.Priority = input.Priority
	}
	if input.BoardID != 0 {
		task.BoardID = input.BoardID
	}

	if input.Deadline != nil {
		if *input.Deadline == "" {
			task.Deadline = nil
		} else {
			t, err := time.Parse(time.RFC3339, *input.Deadline)
			if err != nil {
				t, err = time.Parse("2006-01-02", *input.Deadline)
			}
			if err == nil {
				task.Deadline = &t
			}
		}
	}

	db.DB.Save(&task)

	go func() {
		SendNotification("Task updated: " + task.Title)
	}()

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var task models.Task

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	db.DB.Delete(&task)

	go func() {
		SendNotification("Task deleted: " + task.Title)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func GetTasksByBoard(c *gin.Context) {
	boardID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var tasks []models.Task

	query := db.DB.Where("board_id = ?", boardID)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	query.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	var task models.Task

	query := db.DB.Where("id = ?", id)
	if role != "admin" && role != "moderator" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{"todo": true, "in-progress": true, "done": true}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use: todo, in-progress, done"})
		return
	}

	task.Status = input.Status
	db.DB.Save(&task)

	go func() {
		SendNotification("Task status changed to: " + task.Status)
	}()

	c.JSON(http.StatusOK, task)
}
