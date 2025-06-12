package events

import (
	"fmt"
	"log/slog"
	"pos-master/models"
	"pos-master/proto/dashboard"
	dashboardservices "pos-master/services/dashboard_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
)

func GetEventsHandler(c *gin.Context) {

	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, fmt.Sprintf("error: %v", err))
		return
	}

	getRequest.SetDefaults()

	req := &dashboard.GetEventsRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	events, err := dashboardservices.GetEvents(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve events ", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("events"), fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Events"), gin.H{
		"events": events,
	})

}
