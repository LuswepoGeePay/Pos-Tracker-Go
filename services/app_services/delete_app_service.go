package appservices

import (
	"fmt"
	"log/slog"
	database "pos-master/config"
	"pos-master/models"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"

	"github.com/google/uuid"
)

func DeleteApp(appID string) error {

	parsedAppID, err := uuid.Parse(appID)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse app id", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to parse app id %v", err))
	}

	tx := database.DB.Begin()

	res1 := tx.Unscoped().Delete(&models.AppVersion{}, "app_id = ?", parsedAppID)
	if err := res1.Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete app version", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete app version %v", err))

	}
	utils.Log(slog.LevelInfo, "✅info", "deleted app versions", "count", res1.RowsAffected)

	res2 := tx.Unscoped().Delete(&models.App{}, "id = ?", parsedAppID)
	if err := res2.Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete app", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete app %v", err))

	}
	utils.Log(slog.LevelInfo, "✅info", "deleted app", "count", res2.RowsAffected)

	if err := tx.Commit().Error; err != nil {
		return utils.CapitalizeError(fmt.Sprintf("unable to commit transaction %v", err))
	}

	eventservices.RegisterEvent("An app has been deleted", map[string]interface{}{
		"App ID": appID,
	})

	return nil
}

func DeleteAppVersion(versionID string) error {
	parsedVersionID, err := uuid.Parse(versionID)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse app id", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to parse app id %v", err))
	}

	tx := database.DB.Begin()
	if err := tx.Unscoped().Delete(&models.AppVersion{}, "id = ?", parsedVersionID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete app version", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete app version %v", err))

	}
	if err := tx.Commit().Error; err != nil {
		return utils.CapitalizeError(fmt.Sprintf("unable to commit transaction %v", err))
	}

	eventservices.RegisterEvent("An app version has been deleted", map[string]interface{}{
		"App version ID": versionID,
	})

	return nil
}
