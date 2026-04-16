package terminaltype

import (
	"fmt"
	"log/slog"
	"strconv"

	"pos-master/models"
	"pos-master/proto/terminaltype"
	terminaltypeservices "pos-master/services/terminal_type_services"
	"pos-master/utils"
	"pos-master/utils/sentry"

	"github.com/gin-gonic/gin"
)

func CreateTerminalTypeHandler(c *gin.Context) {

	var req terminaltype.RegisterTerminalTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	err := terminaltypeservices.CreateTerminalType(&req)
	if err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, "Terminal Type created successfully")

}

func GetTerminalTypesHandler(c *gin.Context) {
	// Read pagination and search parameters from query string
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	searchQuery := c.DefaultQuery("search", "")

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

	req := &terminaltype.GetTerminalTypesRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	terminalTypes, err := terminaltypeservices.GetTerminalTypes(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve terminal types ", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("terminal types"), fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("terminal types"), gin.H{
		"data": terminalTypes,
	})

}

func EditTerminalTypeHandler(c *gin.Context) {

	var req terminaltype.EditTerminalTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	err := terminaltypeservices.EditTerminalType(&req)
	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, "Terminal Type edited successfully")

}

func DeleteTerminalTypeHandler(c *gin.Context) {

	terminalTypeID := c.Param("id")

	err := terminaltypeservices.DeleteTerminalType(terminalTypeID)

	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error : %v", err))
		return
	}

	utils.RespondWithSuccess(c, "POS Device deleted successfully")

}
