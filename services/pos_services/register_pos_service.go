package posservices

import (
	"errors"
	"fmt"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RegisterPosDevice(req *posdevices.RegisterPosDeviceRequest) (string, error) {

	posDeviceID := uuid.New()

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Error("unable to start pos register transaction", "error", tx.Error.Error())
		return "", utils.CapitalizeError("Unable to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var business models.Business
	result := tx.Where("email = ?", req.Email).First(&business)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.Error("unable to find business with provided email",
				"email", req.Email,
				"device_model", req.DeviceModel,
			)
			return "", utils.CapitalizeError("Business not found with provided email")
		}
		utils.Error("unable to find business",
			"email", req.Email,
			"error", result.Error.Error(),
		)
		return "", utils.CapitalizeError(utils.FormatError("unable to find business", result.Error))
	}

	terminalTypeID := resolveTerminalTypeID(tx, req)

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
		TerminalTypeID:             terminalTypeID,
		DeviceIdentificationNumber: req.DeviceIdentificationNumber,
	}

	result = tx.Create(&pos)
	if result.Error != nil {
		tx.Rollback()
		utils.Error("unable to register pos device",
			"email", req.Email,
			"error", result.Error.Error(),
		)
		return "", utils.CapitalizeError("Unable to register pos device")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Error("failed to commit pos register transaction", "error", err.Error())
		return "", utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	terminalTypeValue := ""
	if terminalTypeID != nil {
		terminalTypeValue = terminalTypeID.String()
	}

	utils.Info("pos device registered",
		"device_id", posDeviceID.String(),
		"email", req.Email,
		"business_id", business.ID.String(),
		"device_model", req.DeviceModel,
		"terminal_type_id", terminalTypeValue,
	)

	eventservices.RegisterEvent("POS device registered", map[string]interface{}{
		"Pos ID":           posDeviceID.String(),
		"Serial number":    req.SerialNumber,
		"Business Name":    business.Name,
		"Description":      req.Description,
		"Device Model":     req.DeviceModel,
		"Status":           "online",
		"Operating system": req.OperatingSystem,
		"Terminal type":    terminalTypeValue,
	})

	return posDeviceID.String(), nil
}

func resolveTerminalTypeID(tx *gorm.DB, req *posdevices.RegisterPosDeviceRequest) *uuid.UUID {
	rawID := strings.TrimSpace(req.GetTerminalTypeId())
	if rawID != "" {
		parsed, err := uuid.Parse(rawID)
		if err == nil {
			var terminalType models.TerminalType
			if err := tx.Where("id = ?", parsed).First(&terminalType).Error; err == nil {
				utils.Info("resolved terminal type from id", "terminal_type_id", parsed.String())
				return &parsed
			}
			utils.Warn("terminal type id not found, trying device model",
				"terminal_type_id", rawID,
				"device_model", req.GetDeviceModel(),
			)
		} else {
			utils.Warn("terminal_type_id is not a UUID, trying device model",
				"terminal_type_id", rawID,
				"device_model", req.GetDeviceModel(),
			)
		}
	}

	model := strings.TrimSpace(req.GetDeviceModel())
	if model == "" {
		utils.Info("no terminal type id or device model provided, storing null terminal type",
			"email", req.GetEmail(),
		)
		return nil
	}

	var terminalType models.TerminalType
	err := tx.Where("LOWER(terminal_model) = LOWER(?) OR LOWER(name) = LOWER(?)", model, model).
		First(&terminalType).Error
	if err != nil {
		utils.Warn("could not resolve terminal type from device model, storing null",
			"device_model", model,
			"error", err.Error(),
		)
		return nil
	}

	utils.Info("resolved terminal type from device model",
		"device_model", model,
		"terminal_type_id", terminalType.ID.String(),
		"terminal_type_name", terminalType.Name,
	)
	return &terminalType.ID
}
