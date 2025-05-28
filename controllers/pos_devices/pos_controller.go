package posdevices

import (
	"pos-master/proto/posdevices"
	posservices "pos-master/services/pos_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func RegisterPosDeviceHandler(c *gin.Context) {

	var req posdevices.RegisterPosDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, err.Error())
		return
	}

	err := posservices.RegisterPosDevice(&req)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error())
	}

	utils.RespondWithSuccess(c, "POS Registered successfully")

}
