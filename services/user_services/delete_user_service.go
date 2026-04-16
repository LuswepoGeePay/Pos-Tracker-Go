package userservices

import (
	database "pos-master/config"
	"pos-master/models"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"

	"github.com/google/uuid"
)

func DeleteUser(userId string) error {
	parsedID, err := uuid.Parse(userId)
	if err != nil {
		return utils.CapitalizeError("invalid ID format")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Unscoped().Delete(&models.User{}, "id = ?", parsedID).Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError("failed to delete user")
	}

	tx.Commit()
	eventservices.RegisterEvent("User deleted successfully", map[string]interface{}{
		"User ID": userId,
	})
	return nil
}
