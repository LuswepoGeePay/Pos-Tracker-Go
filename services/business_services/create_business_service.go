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

func CreateBusiness(c *gin.Context, req *business.BusinessRegisterRequest) error {

	token, err := pocketbase.HandlePocketBaseAuth(c)

	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to get pocketbase token", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to get pocketbase token")
	}

	fileURL, err := pocketbase.HandleUpload(c, token, "file")
	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to upload file to pocketbase", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to upload file to server")
	}

	newBusiness := models.Business{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		Address:      req.Address,
		Status:       true,
		Phone:        req.Phone,
		BusinessLogo: fileURL,
	}

	// Start transaction for business creation
	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to start transaction: %v", tx.Error))
	}

	result := tx.Create(&newBusiness)
	if result.Error != nil {
		tx.Rollback()
		return utils.CapitalizeError(utils.FormatError("unable to create business", result.Error))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	eventservices.RegisterEvent("New business Registered", map[string]interface{}{
		"name":    req.Name,
		"address": req.Address,
		"email":   req.Email,
		"phone":   req.Phone,
	})

	return nil
}
