package auth

import (
	"fmt"
	pb "pos-master/proto/auth"
	services "pos-master/services/authservices"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func UpdatePasswordHandler(c *gin.Context) {

	var req pb.ResetPasswordRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, 400, utils.FailBind)
		return
	}

	err = services.ResetPassword(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Failed to reset password.", fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, "Password has been reset successfully")

}
