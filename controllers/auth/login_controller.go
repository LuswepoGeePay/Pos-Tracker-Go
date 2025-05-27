package auth

import (
	"net/http"
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	services "pos-master/services/authservices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoginHandler(c *gin.Context) {
	var req pb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "Invalid request format",
		})
		return
	}

	response, err := services.LoginUser(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failure",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
func AssignPermissionToRole(roleID uuid.UUID, permissionID uuid.UUID) error {
	role := models.Role{}
	permission := models.Permission{}

	// Get the role and permission from the database
	config.DB.First(&role, "id = ?", roleID)
	config.DB.First(&permission, "id = ?", permissionID)

	// Associate the permission with the role
	config.DB.Model(&role).Association("Permissions").Append(&permission)

	return nil
}
