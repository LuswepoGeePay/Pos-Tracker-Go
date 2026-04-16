package posservices

import (
	"errors"
	"fmt"
	"log/slog"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RegisterPosDevice(req *posdevices.RegisterPosDeviceRequest) (string, error) {

	posDeviceID := uuid.New()

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var business models.Business
	result := tx.Where("email = ?", req.Email).First(&business)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.Log(slog.LevelError, "❌error", "unable to find business with provided email", "data", req)

			return "", utils.CapitalizeError("Business not found with provided email")
		}
		utils.Log(slog.LevelError, "❌error", "unable to find business", "deatil", result.Error, "data", req)

		return "", utils.CapitalizeError(utils.FormatError("unable to find business", result.Error))
	}

	terminalTypeID, err := uuid.Parse(req.TerminalTypeId)
	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse terminal type id", "data", req)
		return "", utils.CapitalizeError("Unable to parse terminal type id")
	}

	var terminalType models.TerminalType
	result = database.DB.Where("id = ?", terminalTypeID).First(&terminalType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.Log(slog.LevelError, "❌error", "unable to find terminal type with provided id", "data", req)

			return "", utils.CapitalizeError("Terminal type not found with provided id")
		}
		utils.Log(slog.LevelError, "❌error", "unable to find terminal type", "deatil", result.Error, "data", req)

		return "", utils.CapitalizeError(utils.FormatError("unable to find terminal type", result.Error))
	}

	pos := models.PosDevice{
		ID:                         posDeviceID,
		SerialNumber:               req.SerialNumber,
		Name:                       req.Name,
		Description:                req.Description,
		CurrentAppVersion:          req.CurrentAppVersion,
		LastKnownLatitude:          req.LastKnownLatitude,
		LastKnownLongitude:         req.LastKnownLongitude,
		DeviceModel:                req.DeviceModel,
		OperatingSystem:            req.OperatingSystem,
		Status:                     "online",
		LocationLastUpdatedAt:      time.Now(),
		Email:                      req.Email,
		Entity:                     business.Name,
		FingerPrint:                req.Fingerprint,
		BusinessID:                 business.ID,
		PhoneNumber1:               req.PrimaryNumber,
		PhoneNumber2:               req.SecondaryNumber,
		TerminalTypeID:             &terminalTypeID,
		DeviceIdentificationNumber: req.DeviceIdentificationNumber,
	}

	result = tx.Create(&pos)
	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error", "Unable to register pos device", "Detail", result.Error.Error(), "data", req)
		return "", utils.CapitalizeError("Unable to register pos device")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error", "Failed to commit transaction", "Detail", err.Error())
		return "", utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}
	eventservices.RegisterEvent("POS device registered", map[string]interface{}{
		"Pos ID":           pos,
		"Serial number":    req.SerialNumber,
		"Business Name":    req.BusinessName,
		"Description":      req.Description,
		"Device Model":     req.DeviceModel,
		"Status":           req.Status,
		"Operating system": req.OperatingSystem,
	})

	return posDeviceID.String(), nil

}
