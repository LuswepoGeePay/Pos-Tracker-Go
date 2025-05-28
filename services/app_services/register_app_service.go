package appservices

import (
	"fmt"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	"pos-master/services/pocketbase"
	"pos-master/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterApp(req *appPb.RegisterAppRequest) error {

	app := models.App{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	result := config.DB.Create(&app)
	if result.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to register app %s", result.Error.Error()))
	}

	return nil

}

func RegisterAppVersion(c *gin.Context, req *appPb.RegisterAppVersionRequest) error {

	token, err := pocketbase.HandlePocketBaseAuth(c)

	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to get pocketbase token", "detail", err.Error())
		return utils.CapitalizeError("unable to get pocketbase token")
	}

	fileURL, err := pocketbase.HandleUpload(c, token)

	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to upload file to pocketbase", err.Error())
		return utils.CapitalizeError("unable to upload file to server")
	}

	appID, err := uuid.Parse(req.AppId)

	if err != nil {
		utils.Log(slog.LevelError, "error", err.Error())
		return utils.CapitalizeError("")
	}

	checkSum := ""

	appVersion := models.AppVersion{
		ID:             uuid.New(),
		AppID:          appID,
		BuildNumber:    req.BuildNumber,
		ReleaseNotes:   req.ReleaseNotes,
		FilePath:       fileURL,
		FileSizeBytes:  req.FileSizeBytes,
		CheckSum:       checkSum,
		IsActive:       false,
		IsLatestStable: req.IsLatestStable,
		ReleasedAt:     time.Now(),
		VersionNumber:  req.VersionNumber,
	}

	result := config.DB.Create(&appVersion)

	if result.Error != nil {

	}

	return nil

}
