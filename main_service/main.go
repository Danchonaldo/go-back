package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-proj/main_service/handlers"
)

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {})
	r.POST("/task", handlers.CreateTask)
	r.GET("/tasks", handlers.GetTasks)
	r.Run(":8082")
}
