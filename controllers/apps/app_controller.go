package apps

import (
	"log/slog"
	"pos-master/models"
	appPb "pos-master/proto/app"
	appservices "pos-master/services/app_services"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

func RegisterAppHandler(c *gin.Context) {

	var req appPb.RegisterAppRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, err.Error())
		return
	}

	err := appservices.RegisterApp(&req)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error())
	}

	utils.RespondWithSuccess(c, "App Registered successfully")

}

func RegisterNewAppVersionHandler(c *gin.Context) {

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.Log(slog.LevelError, "❌error", "failed to parse form")
		utils.RespondWithError(c, 400, "Failed to parse form", err.Error())
		return
	}

	versionData := c.Request.FormValue("version")

	if versionData == "" {
		utils.Log(slog.LevelError, "❌error", "invalid app version data")
		utils.RespondWithError(c, 400, "App version Data is missing")
		return
	}

	var req appPb.RegisterAppVersionRequest

	if err := protojson.Unmarshal([]byte(versionData), &req); err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to marshal app data")
		utils.RespondWithError(c, 400, "Unable to marshal app data", err.Error())
		return
	}

	err := appservices.RegisterAppVersion(c, &req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", err.Error())
		utils.RespondWithError(c, 400, "Unable to upload new app to version", err.Error())
		return

	}

	utils.RespondWithSuccess(c, "Added new app version data")

}

func GetAppsHandler(c *gin.Context) {

	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &appPb.GetAppsRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	apps, err := appservices.GetApps(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve apps", "details", string(err.Error()))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Apps"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Apps"), gin.H{
		"apps": apps,
	})

}

func GetAppVersionsHandler(c *gin.Context) {
	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &appPb.GetAppVersionsRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	apps, err := appservices.GetAppVersions(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve app versions ", "details", string(err.Error()))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Apps versions"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Apps versions"), gin.H{
		"apps": apps,
	})

}
