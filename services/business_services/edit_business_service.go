package businessservices

import (
	"fmt"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/business"
	eventservices "pos-master/services/event_services"
	"pos-master/services/pocketbase"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func EditBusiness(c *gin.Context, req *business.EditBusinessRequest) error {
	businessID, err := uuid.Parse(req.Id)
	if err != nil {
		utils.Error("failed to parse business ID", "id", req.Id, "error", err.Error())
		return utils.CapitalizeError(fmt.Sprintf("failed to parse business ID %v", err))
	}

	var currentBusiness models.Business
	result := database.DB.Where("id = ?", businessID).First(&currentBusiness)
	if result.Error != nil {
		utils.Error("unable to find business for edit", "id", req.Id, "error", result.Error.Error())
		return utils.CapitalizeError("unable to find business with that ID")
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	if req.Address != "" {
		updates["address"] = req.Address
	}

	_, err = c.FormFile("file")
	if err == nil {
		token, err := pocketbase.HandlePocketBaseAuth(c)
		if err != nil {
			utils.Error("unable to get pocketbase token", "error", err.Error(), "business_id", req.Id)
			return utils.CapitalizeError("unable to get pocketbase token")
		}

		fileUrl, err := pocketbase.HandleUpload(c, token, "file")
		if err != nil {
			utils.Error("unable to upload file to pocketbase", "error", err.Error(), "business_id", req.Id)
			return utils.CapitalizeError("unable to upload file to server")
		}

		updates["business_logo"] = fileUrl
	}

	if req.Status != currentBusiness.Status {
		updates["status"] = req.Status
	}

	if req.Phone != currentBusiness.Phone {
		updates["phone"] = req.Phone
	}

	if len(updates) == 0 {
		return utils.CapitalizeError("no changes detected")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to start transaction: %v", tx.Error))
	}

	err = tx.Model(&models.Business{}).Where("id = ?", businessID).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		utils.Error("failed to update business", "business_id", req.Id, "error", err.Error())
		return utils.CapitalizeError(fmt.Sprintf("failed to update business: %v", err))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	utils.Info("business updated", "business_id", req.Id)

	eventservices.RegisterEvent("business edited", map[string]interface{}{
		"name":    req.Name,
		"address": req.Address,
		"phone":   req.Phone,
		"email":   req.Email,
	})

	return nil
}
