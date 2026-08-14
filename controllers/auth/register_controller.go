package auth

import (
	"fmt"
	pb "pos-master/proto/auth"
	userservices "pos-master/services/user_services"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var req pb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "failure",
			"error":  fmt.Sprintf("error: %v", err),
		})
		return
	}

	response := userservices.RegisterUser(&req)
	if !response.Success {
		c.JSON(400, gin.H{
			"status":  "failure",
			"error":   response.Message,
			"success": false,
			"message": response.Message,
		})
		return
	}

	c.JSON(200, response)
}
