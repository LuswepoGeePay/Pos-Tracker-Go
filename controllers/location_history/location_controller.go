package locationhistory

import (
	"pos-master/proto/posdevices"
	historyservices "pos-master/services/history_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func RegisterNewLocationHandler(c *gin.Context) {

	var req posdevices.RegisterLocationHistoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, err.Error())
		return
	}

	err := historyservices.RegisterNewLocationHistory(&req)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error())
	}

	utils.RespondWithSuccess(c, "Location Registered successfully")

}
