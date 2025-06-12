package posdevices

import (
	"fmt"
	"log/slog"
	"pos-master/models"
	"pos-master/proto/posdevices"
	posservices "pos-master/services/pos_services"
	"pos-master/utils"
	"pos-master/utils/sentry"

	"github.com/gin-gonic/gin"
)

func RegisterPosDeviceHandler(c *gin.Context) {

	var req posdevices.RegisterPosDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	posDeviceID, err := posservices.RegisterPosDevice(&req)
	if err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
	}

	utils.RespondWithSuccess(c, "POS Registered successfully", gin.H{
		"device_id": posDeviceID,
	})

}

func GetPosDevicesHandler(c *gin.Context) {
	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, fmt.Sprintf("error: %v", err))
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
		utils.Log(slog.LevelError, "❌error", "unable to retrieve pos devices ", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("pos devices"), fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("pos devices"), gin.H{
		"devices": devices,
	})

}

func EditDeviceHandler(c *gin.Context) {

	var req posdevices.EditPosDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	err := posservices.EditDevice(&req)
	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, "POS Edited successfully")

}

func DeleteDeviceHandler(c *gin.Context) {

	deviceID := c.Param("id")

	err := posservices.DeleteDevice(deviceID)

	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error : %v", err))
		return
	}

	utils.RespondWithSuccess(c, "POS Device deleted successfully")

}

func GetDeviceByID(c *gin.Context) {

}
