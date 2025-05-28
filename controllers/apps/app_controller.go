package apps

import (
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

func RegisterNewAppVersion(c *gin.Context) {

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondWithError(c, 400, "Failed to parse form", err.Error())
		return
	}

	versionData := c.Request.FormValue("version")

	if versionData == "" {
		utils.RespondWithError(c, 400, "App version Data is missing")
		return
	}

	var req appPb.RegisterAppVersionRequest

	if err := protojson.Unmarshal([]byte(versionData), &req); err != nil {
		utils.RespondWithError(c, 400, "Invalid even data", err.Error())
		return
	}

}

func GetAppsHandler(c *gin.Context) {

}
