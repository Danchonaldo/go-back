package handlers

import (
	"go-proj/main_service/db"
	"go-proj/main_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	db.DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	c.JSON(http.StatusOK, users)
}
