package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

var notifications []Notification
var counter int

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "notification_service",
			"time":    time.Now(),
		})
	})

	r.POST("/notify", func(c *gin.Context) {
		var req map[string]string
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		message, ok := req["message"]
		if !ok || message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
			return
		}

		counter++
		notif := Notification{
			ID:        counter,
			Message:   message,
			CreatedAt: time.Now(),
		}
		notifications = append(notifications, notif)

		log.Printf("[NOTIFICATION] %s", message)

		c.JSON(http.StatusOK, gin.H{
			"status":       "notification sent",
			"notification": notif,
		})
	})

	r.GET("/notifications", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"count":         len(notifications),
			"notifications": notifications,
		})
	})

	log.Println("Notification service running on :8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("Failed to start notification service:", err)
	}
}
