package posdevices

import (
	"fmt"
	"log/slog"
	"strconv"

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
		return
	}

	utils.RespondWithSuccess(c, "POS Registered successfully", gin.H{
		"device_id": posDeviceID,
	})

}

func GetPosDevicesHandler(c *gin.Context) {
	// Read pagination and search parameters from query string
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	searchQuery := c.DefaultQuery("search", "")

	// Read filter parameters
	businessId := c.DefaultQuery("business_id", "")
	status := c.DefaultQuery("status", "")
	phoneNumber := c.DefaultQuery("phone_number", "")
	appVersion := c.DefaultQuery("current_app_version", "")
	serialNumber := c.DefaultQuery("serial_number", "")
	locationStartDate := c.DefaultQuery("location_last_updated_start", "")
	locationEndDate := c.DefaultQuery("location_last_updated_end", "")

	// Parse page and pageSize to integers
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1
	}
	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeNum = 10
	}

	getRequest := models.SearchRequest{
		GetRequest: models.GetRequest{
			Page:     pageNum,
			PageSize: pageSizeNum,
		},
		SearchQuery: searchQuery,
	}

	req := &posdevices.GetPosDevicesRequest{
		Page:                     int32(getRequest.Page),
		PageSize:                 int32(getRequest.PageSize),
		SearchQuery:              getRequest.SearchQuery,
		BusinessId:               businessId,
		Status:                   status,
		PhoneNumber:              phoneNumber,
		AppVersion:               appVersion,
		SerialNumber:             serialNumber,
		LocationLastUpdatedStart: locationStartDate,
		LocationLastUpdatedEnd:   locationEndDate,
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

func HeartBeatHandler(c *gin.Context) {

	var req models.HeartBeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to bind request", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, "error: invalid payload")
		return
	}

	utils.Log(slog.LevelInfo, "✅info", "heartbeat", "details", string(fmt.Sprintf("request: %v", req)))

	err := posservices.RegisterHeartBeat(&req)
	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to register heartbeat", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, "error", utils.FormatError("detail", err))
		return
	}

	utils.RespondWithSuccess(c, "heartbeat recieved")
}
