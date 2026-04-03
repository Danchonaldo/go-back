package main

import (
	"go-proj/db"
	"go-proj/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db.ConnectDB()

	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)

	r.POST("/boards", handlers.CreateBoard)
	r.GET("/boards", handlers.GetBoards)

	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks", handlers.GetTasks)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.GET("/boards/:id/tasks", handlers.GetTasksByBoard)

	r.Run(":8080")
}
