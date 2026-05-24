package main

import (
	"log"
	"main_service/db"
	"main_service/handlers"
	"main_service/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
	db.RunMigrations()

	// Создаём встроенного админа
	handlers.SeedAdminUser()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Logger
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%s] %s %s - %d (%v)",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "main_service"})
	})

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		users := api.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUserByID)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		boards := api.Group("/boards")
		{
			boards.POST("", handlers.CreateBoard)
			boards.GET("", handlers.GetBoards)
			boards.GET("/:id", handlers.GetBoardByID)
			boards.PUT("/:id", handlers.UpdateBoard)
			boards.DELETE("/:id", handlers.DeleteBoard)
			boards.GET("/:id/tasks", handlers.GetTasksByBoard)
		}

		tasks := api.Group("/tasks")
		{
			tasks.POST("", handlers.CreateTask)
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTaskByID)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
			tasks.PATCH("/:id/status", handlers.UpdateTaskStatus)
		}

		// Admin only routes
		admin := api.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/users", handlers.AdminGetUsers)
			admin.PUT("/users/:id/role", handlers.AdminUpdateUserRole)
			admin.DELETE("/users/:id", handlers.AdminDeleteUser)
			admin.GET("/stats", handlers.AdminGetStats)
			admin.GET("/boards", handlers.AdminGetAllBoards)
			admin.GET("/tasks", handlers.AdminGetAllTasks)
		}
	}

	r.POST("/notify", func(c *gin.Context) {
		result, err := handlers.SendNotification("manual notification")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Notification service unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	log.Println("Main service running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
