package handlers

import (
	"main_service/db"
	"main_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminGetUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	result := make([]gin.H, 0, len(users))
	for _, u := range users {
		result = append(result, gin.H{
			"id":         u.ID,
			"name":       u.Name,
			"email":      u.Email,
			"role":       u.Role,
			"created_at": u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, result)
}

func AdminUpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	currentUserID, _ := c.Get("userID")

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ID == currentUserID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change your own role"})
		return
	}

	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Только user и admin
	validRoles := map[string]bool{"user": true, "admin": true}
	if !validRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Use: user, admin"})
		return
	}

	user.Role = input.Role
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	currentUserID, _ := c.Get("userID")

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ID == currentUserID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	db.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func AdminGetAllBoards(c *gin.Context) {
	var boards []models.Board
	db.DB.Find(&boards)
	c.JSON(http.StatusOK, boards)
}

func AdminGetAllTasks(c *gin.Context) {
	var tasks []models.Task
	db.DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func AdminGetStats(c *gin.Context) {
	var userCount, boardCount, taskCount int64
	var todoCount, inProgressCount, doneCount int64

	db.DB.Model(&models.User{}).Count(&userCount)
	db.DB.Model(&models.Board{}).Count(&boardCount)
	db.DB.Model(&models.Task{}).Count(&taskCount)
	db.DB.Model(&models.Task{}).Where("status = ?", "todo").Count(&todoCount)
	db.DB.Model(&models.Task{}).Where("status = ?", "in-progress").Count(&inProgressCount)
	db.DB.Model(&models.Task{}).Where("status = ?", "done").Count(&doneCount)

	c.JSON(http.StatusOK, gin.H{
		"users":       userCount,
		"boards":      boardCount,
		"tasks":       taskCount,
		"todo":        todoCount,
		"in_progress": inProgressCount,
		"done":        doneCount,
	})
}
