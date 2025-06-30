package users

import (
	"fmt"
	"log/slog"
	"pos-master/models"
	"pos-master/proto/auth"
	userservices "pos-master/services/user_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func GetUsersHandler(c *gin.Context) {

	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, fmt.Sprintf("error: %v", err))
		return
	}

	getRequest.SetDefaults()

	req := &auth.GetUsersRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	users, err := userservices.GetUsers(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve users", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Users"), fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Users"), gin.H{
		"users": users,
	})

}

func EditUserHandler(c *gin.Context) {
	var req auth.EditUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "failure",
			"message": "Invalid request format",
		})
		return
	}

	err := userservices.EditUser(&req)
	if err != nil {
		utils.RespondWithError(c, 400, utils.FormatError("unable to update user", err))
		return
	}

	utils.RespondWithSuccess(c, "User updated successfully")

}

func GetUserHandler(c *gin.Context) {

	userid := c.Param("user_id")

	user, err := userservices.GetUser(userid)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FormatError("unable to get user info", err))
		return
	}

	utils.RespondWithSuccess(c, "✅ user info retrieved", gin.H{
		"data": user,
	})

}
