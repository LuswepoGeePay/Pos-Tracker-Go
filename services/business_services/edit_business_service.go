package businessservices

import (
	"fmt"
	"log/slog"
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
		utils.Log(slog.LevelError, "error", "failed to parse business ID")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse business ID %v", fmt.Sprintf("error: %v", err)))
	}
	updates := map[string]interface{}{}

	var currentAppVersion models.Business

	result := database.DB.Where("id = ?", businessID).Find(&currentAppVersion)
	if result.Error != nil {

	}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	_, err = c.FormFile("file")

	if err == nil {
		// File was uploaded, handle upload
		token, err := pocketbase.HandlePocketBaseAuth(c)
		if err != nil {
			utils.Log(slog.LevelError, "error", "unable to get pocketbase token", "detail", err.Error())
			return utils.CapitalizeError("unable to get pocketbase token")
		}

		fileUrl, err := pocketbase.HandleUpload(c, token, "file")
		if err != nil {
			utils.Log(slog.LevelError, "error", "unable to upload file to pocketbase", err.Error())
			return utils.CapitalizeError("unable to upload file to server")
		}

		updates["business_logo"] = fileUrl
	}

	if req.Status != currentAppVersion.Status {
		updates["status"] = req.Status
	}

	if req.Phone != currentAppVersion.Phone {
		updates["phone"] = req.Phone
	}

	tx := database.DB.Begin()

	err = tx.Model(&models.AppVersion{}).Where("id = ?", businessID).Updates(updates).Error

	if err != nil {
		return utils.CapitalizeError(fmt.Sprintf("failed to update business: %v", fmt.Sprintf("error: %v", err)))
	}

	eventservices.RegisterEvent("business edited", map[string]interface{}{
		"name":    req.Name,
		"address": req.Address,
		"phone":   req.Phone,
		"email":   req.Email,
	})

	tx.Commit()

	return nil

}
