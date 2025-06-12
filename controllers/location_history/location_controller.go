package locationhistory

import (
	"fmt"
	"log/slog"
	"pos-master/models"
	"pos-master/proto/posdevices"
	historyservices "pos-master/services/history_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func RegisterNewLocationHandler(c *gin.Context) {

	var req posdevices.RegisterLocationHistoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	err := historyservices.RegisterNewLocationHistory(&req)
	if err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, "Location Registered successfully")

}

func GetLocationsHandler(c *gin.Context) {
	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, fmt.Sprintf("error: %v", err))
		return
	}

	getRequest.SetDefaults()

	req := &posdevices.GetLocationHistorysRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	history, err := historyservices.GetLocationHistory(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve location history ", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("location history"), fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("location history"), gin.H{
		"history": history,
	})

}
