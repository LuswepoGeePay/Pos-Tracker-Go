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

func CreateBusiness(c *gin.Context, req *business.BusinessRegisterRequest) error {
	utils.Info("create business attempt",
		"name", req.Name,
		"email", req.Email,
		"phone", req.Phone,
	)

	fileURL := ""
	_, fileErr := c.FormFile("file")
	if fileErr == nil {
		token, err := pocketbase.HandlePocketBaseAuth(c)
		if err != nil {
			utils.Error("unable to get pocketbase token", "error", err.Error(), "email", req.Email)
			return utils.CapitalizeError("unable to get pocketbase token")
		}

		fileURL, err = pocketbase.HandleUpload(c, token, "file")
		if err != nil {
			utils.Error("unable to upload file to pocketbase", "error", err.Error(), "email", req.Email)
			return utils.CapitalizeError("unable to upload file to server")
		}
	} else {
		utils.Info("creating business without logo", "email", req.Email)
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

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Error("unable to start business create transaction", "error", tx.Error.Error())
		return utils.CapitalizeError(fmt.Sprintf("Unable to start transaction: %v", tx.Error))
	}

	result := tx.Create(&newBusiness)
	if result.Error != nil {
		tx.Rollback()
		utils.Error("unable to create business", "email", req.Email, "error", result.Error.Error())
		return utils.CapitalizeError(utils.FormatError("unable to create business", result.Error))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Error("failed to commit business create transaction", "error", err.Error())
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	utils.Info("business created",
		"business_id", newBusiness.ID.String(),
		"email", req.Email,
		"has_logo", fileURL != "",
	)

	eventservices.RegisterEvent("New business Registered", map[string]interface{}{
		"name":    req.Name,
		"address": req.Address,
		"email":   req.Email,
		"phone":   req.Phone,
	})

	return nil
}
