package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-proj/main_service/db"
	"go-proj/main_service/handlers"
	"go-proj/main_service/middleware"
)

func main() {
	db.ConnectDB()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/board", handlers.CreateBoard)
		auth.GET("/boards", handlers.GetBoards)

		auth.POST("/task", handlers.CreateTask)
		auth.GET("/tasks", handlers.GetTasks)
		auth.GET("/task/:id", handlers.GetTaskByID)
		auth.PUT("/task/:id", handlers.UpdateTask)
		auth.DELETE("/task/:id", handlers.DeleteTask)
		auth.GET("/board/:id/tasks", handlers.GetTasksByBoard)
	}

	r.Run(":8080")
}
