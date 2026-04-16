package appservices

import (
	"fmt"
	"log/slog"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
)

func CheckAppUpdate(req *posdevices.CheckUpdateRequest) (*posdevices.CheckUpdateResponse, error) {

	var latestVersion models.AppVersion
	var posDevice models.PosDevice

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Fetch POS device first to get its terminal type
	utils.Log(slog.LevelInfo, "✅info", "check app update", "details", fmt.Sprintf("request: %v", req))
	err := tx.Where("id = ?", req.PosdeviceId).First(&posDevice).Error

	if err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to find pos device with this ID", "details", fmt.Sprintf("error: %v", err))
		return nil, utils.CapitalizeError(utils.FormatError("unable to find pos device with this ID", err))
	}

	// 2. Fetch latest app version specific to this terminal type
	// This will look for active, latest stable versions matching the device's TerminalTypeID
	err = tx.Where("is_active = ? AND is_latest_stable = ? AND terminal_type_id = ?", true, true, posDevice.TerminalTypeID).
		Order("released_at desc").
		First(&latestVersion).Error

	if err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "could not fetch latest app version", "details", fmt.Sprintf("terminal_type_id: %v, error: %v", posDevice.TerminalTypeID, err))
		return nil, utils.CapitalizeError("could not fetch latest app version for this terminal type")
	}

	// 3. Update device record's current app version if it doesn't match the specific latest
	if posDevice.CurrentAppVersion != latestVersion.VersionNumber {
		result := tx.Model(&models.PosDevice{}).Where("id = ?", req.PosdeviceId).Update("current_app_version", latestVersion.VersionNumber)
		if result.Error != nil {
			tx.Rollback()
			utils.Log(slog.LevelError, "❌error", "unable to update pos device app version", "details", fmt.Sprintf("error: %v", result.Error))
			return nil, utils.CapitalizeError(utils.FormatError("unable to update pos device app version", result.Error))
		}
	}

	// Commit the transaction to save the updated app version status on the device
	if err := tx.Commit().Error; err != nil {
		return nil, utils.CapitalizeError("failed to commit update transaction")
	}

	// 4. Check if the device already has the latest version matching its type
	if latestVersion.VersionNumber == req.AppVersion {
		return &posdevices.CheckUpdateResponse{
			UpdateAvailable: false,
			LatestVersion:   latestVersion.VersionNumber,
			ReleaseNotes:    latestVersion.ReleaseNotes,
			DownloadUrl:     "",
			Code:            0,
		}, nil
	}

	// 5. Register the update check event
	eventservices.RegisterEvent("Pos Device checked for an update", map[string]interface{}{
		"pos_device_id":     req.PosdeviceId,
		"requested_version": req.AppVersion,
		"latest_version":    latestVersion.VersionNumber,
		"terminal_type":     posDevice.TerminalTypeID,
	})

	return &posdevices.CheckUpdateResponse{
		UpdateAvailable: true,
		LatestVersion:   latestVersion.VersionNumber,
		ReleaseNotes:    latestVersion.ReleaseNotes,
		DownloadUrl:     latestVersion.FilePath,
		Code:            0,
	}, nil
}
