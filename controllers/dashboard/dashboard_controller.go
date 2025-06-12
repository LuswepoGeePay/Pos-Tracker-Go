package dashboard

import (
	dashboardservices "pos-master/services/dashboard_services"
	"pos-master/utils"
	"pos-master/utils/sentry"

	"github.com/gin-gonic/gin"
)

func GetTileInfoHandler(c *gin.Context) {

	response, err := dashboardservices.GetTileInfo()

	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, utils.FormatError("unable to get info", err))
		return
	}

	utils.RespondWithSuccess(c, "info fetched", gin.H{
		"info": response,
	})

}

func GetRecentEventsHandler(c *gin.Context) {

}

func GetPieChartDataHandler(c *gin.Context) {

}

func GetLineChartHandler(c *gin.Context) {

}
