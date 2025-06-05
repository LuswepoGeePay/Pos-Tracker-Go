package posdevices

import (
	"log/slog"
	"pos-master/models"
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

	posDeviceID, err := posservices.RegisterPosDevice(&req)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error())
	}

	utils.RespondWithSuccess(c, "POS Registered successfully", gin.H{
		"device_id": posDeviceID,
	})

}

func GetPosDevicesHandler(c *gin.Context) {
	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &posdevices.GetPosDevicesRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	devices, err := posservices.GetPosDevices(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve pos devices ", "details", string(err.Error()))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("pos devices"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("pos devices"), gin.H{
		"devices": devices,
	})

}
