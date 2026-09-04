package auth

import (
	"fmt"
	pb "pos-master/proto/auth"
	userservices "pos-master/services/user_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	utils.Info("create user request received")

	var req pb.RegisterRequest
	if err := utils.BindProtoJSON(c, &req); err != nil {
		utils.Error("failed to bind create user request", "error", err.Error())
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	utils.Info("create user attempt",
		"email", req.GetEmail(),
		"fullname", req.GetFullname(),
		"role", req.GetRole(),
	)

	response := userservices.RegisterUser(&req)
	if !response.Success {
		utils.Error("create user failed",
			"email", req.GetEmail(),
			"role", req.GetRole(),
			"error", response.Message,
		)
		utils.RespondWithError(c, 400, response.Message)
		return
	}

	c.JSON(200, response)
}
