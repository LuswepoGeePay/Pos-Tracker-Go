package appservices

import (
	"fmt"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/utils"

	"github.com/google/uuid"
)

func DeleteApp(appID string) error {

	parsedAppID, err := uuid.Parse(appID)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse app id", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to parse app id %v", err))
	}

	tx := config.DB.Begin()
	if err := tx.Delete(&models.AppVersion{}, "app_id = ?", parsedAppID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete app", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete app %v", err))

	}
	if err := tx.Delete(&models.App{}, "id = ?", parsedAppID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete app", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete app %v", err))

	}
	tx.Commit()
	return nil
}

func DeleteAppVersion(versionID string) error {
	parsedVersionID, err := uuid.Parse(versionID)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse app id", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to parse app id %v", err))
	}

	tx := config.DB.Begin()
	if err := tx.Unscoped().Delete(&models.AppVersion{}, "id = ?", parsedVersionID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete user", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete user %v", err))

	}
	tx.Commit()

	return nil
}
