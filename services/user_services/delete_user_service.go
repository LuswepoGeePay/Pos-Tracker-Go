package userservices

import (
	"pos-master/config"
	"pos-master/models"
	"pos-master/utils"

	"github.com/google/uuid"
)

func DeleteUser(userId string) error {
	parsedID, err := uuid.Parse(userId)
	if err != nil {
		return utils.CapitalizeError("invalid ID format")
	}

	tx := config.DB.Begin()

	if err := tx.Unscoped().Delete(&models.User{}, "id = ?", parsedID).Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError("failed to delete user")
	}

	tx.Commit()
	return nil
}
