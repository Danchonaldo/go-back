package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/notify", func(c *gin.Context) {
		var req map[string]string

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		c.JSON(200, gin.H{
			"status":  "notification sent",
			"message": req["message"],
		})
	})

	r.Run(":8083")
}
