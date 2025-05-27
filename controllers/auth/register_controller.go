package auth

import (
	pb "pos-master/proto/auth"
	userservices "pos-master/services/user_services"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var req pb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "failure",
			"error":  err.Error(),
		})
		return
	}

	response := userservices.RegisterUser(&req)
	c.JSON(200, response)
}
