package main

import (
	"go-proj/db"
	"go-proj/handlers"
	"go-proj/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db.ConnectDB()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/tasks", handlers.CreateTask)
		auth.GET("/tasks", handlers.GetTasks)
	}

	r.Run(":8080")
}
