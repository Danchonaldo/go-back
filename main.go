package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go-proj/db"
	"go-proj/handlers"
	"go-proj/middleware"
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
