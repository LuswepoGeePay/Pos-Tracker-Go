package businessservices

import (
	"fmt"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"

	"github.com/google/uuid"
)

func DeleteBusiness(businessID string) error {

	parsedBusinessID, err := uuid.Parse(businessID)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse app id", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to parse app id %v", err))
	}

	tx := config.DB.Begin()
	if err := tx.Unscoped().Delete(&models.Business{}, "id = ?", parsedBusinessID).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to delete business", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError(fmt.Sprintf("unable to delete business %v", err))

	}
	tx.Commit()
	eventservices.RegisterEvent("A business has been deleted", map[string]interface{}{
		"business ID": businessID,
	})
	return nil
}
