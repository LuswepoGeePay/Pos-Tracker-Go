package appservices

import (
	"fmt"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	"pos-master/services/pocketbase"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func EditApp(req *appPb.EditAppRequest) error {

	appID, err := uuid.Parse(req.Id)
	if err != nil {
		utils.Log(slog.LevelError, "error", "failed to parse app ID")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse app ID %v", fmt.Sprintf("error: %v", err)))
	}
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	tx := config.DB.Begin()

	err = tx.Model(&models.App{}).Where("id = ?", appID).Updates(updates).Error

	if err != nil {
		return utils.CapitalizeError(fmt.Sprintf("failed to update app: %v", fmt.Sprintf("error: %v", err)))
	}

	tx.Commit()

	return nil
}

func EditAppVersion(c *gin.Context, req *appPb.EditAppVersionRequest) error {

	versionID, err := uuid.Parse(req.VersionId)
	if err != nil {
		utils.Log(slog.LevelError, "error", "failed to parse version ID")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse version ID %v", fmt.Sprintf("error: %v", err)))
	}
	updates := map[string]interface{}{}

	var currentAppVersion models.AppVersion

	result := config.DB.Where("id = ?", versionID).Find(&currentAppVersion)
	if result.Error != nil {

	}

	if req.VersionNumber != "" {
		updates["version_number"] = req.VersionNumber
	}

	if req.ReleaseNotes != "" {
		updates["release_notes"] = req.ReleaseNotes
	}

	_, err = c.FormFile("apk")

	if err == nil {
		// File was uploaded, handle upload
		token, err := pocketbase.HandlePocketBaseAuth(c)
		if err != nil {
			utils.Log(slog.LevelError, "error", "unable to get pocketbase token", "detail", err.Error())
			return utils.CapitalizeError("unable to get pocketbase token")
		}

		apkUrl, err := pocketbase.HandleUpload(c, token, "apk")
		if err != nil {
			utils.Log(slog.LevelError, "error", "unable to upload file to pocketbase", err.Error())
			return utils.CapitalizeError("unable to upload file to server")
		}

		updates["file_path"] = apkUrl
	}

	if req.IsLatestStable != currentAppVersion.IsLatestStable {
		updates["is_latest_stable"] = req.IsLatestStable
	}

	if req.IsActive != currentAppVersion.IsActive {
		updates["is_active"] = req.IsActive
	}

	tx := config.DB.Begin()

	err = tx.Model(&models.App{}).Where("id = ?", versionID).Updates(updates).Error

	if err != nil {
		return utils.CapitalizeError(fmt.Sprintf("failed to update version: %v", fmt.Sprintf("error: %v", err)))
	}

	tx.Commit()

	return nil

}
