package posservices

import (
	"fmt"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"

	"github.com/google/uuid"
)

func DeleteDevice(deviceID string) error {

	parsedID, err := uuid.Parse(deviceID)
	if err != nil {
		utils.Log(slog.LevelError, "error", "invalid device ID")
		return utils.CapitalizeError(fmt.Sprintf("unable to parse pos device ID %v", err))
	}

	tx := config.DB.Begin()

	if err := tx.Delete(&models.LocationHistory{}, "pos_device_id = ?", parsedID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete history", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete  history %v", err))

	}

	if err := tx.Delete(&models.PosDevice{}, "id = ?", parsedID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete device", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete device %v", err))

	}

	tx.Commit()

	eventservices.RegisterEvent("POS Device has been deleted", map[string]interface{}{
		"Pos ID": deviceID})

	return nil
}
