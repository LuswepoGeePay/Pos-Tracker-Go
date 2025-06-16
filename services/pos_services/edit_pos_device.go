package posservices

import (
	"fmt"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"

	"github.com/google/uuid"
)

func EditDevice(req *posdevices.EditPosDeviceRequest) error {
	deviceID, err := uuid.Parse(req.Id)
	if err != nil {
		utils.Log(slog.LevelError, "error", "failed to parse device ID")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse device ID %v", fmt.Sprintf("error: %v", err)))
	}
	updates := map[string]interface{}{}

	var currentDevice models.PosDevice

	result := config.DB.Where("id = ?", deviceID).Find(&currentDevice)
	if result.Error != nil {

	}

	if req.SerialNumber != "" {
		updates["serial_number"] = req.SerialNumber
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.BusinessName != "" {
		updates["entity"] = req.BusinessName
	}

	if req.Description != currentDevice.Description {
		updates["description"] = req.Description
	}

	if req.CurrentAppVersion != currentDevice.CurrentAppVersion {
		updates["current_app_version"] = req.CurrentAppVersion
	}

	if req.LastKnownLatitude != currentDevice.LastKnownLatitude {
		updates["last_known_latitude"] = req.LastKnownLatitude
	}

	if req.LastKnownLongitude != currentDevice.LastKnownLongitude {
		updates["last_known_longitude"] = req.LastKnownLongitude
	}
	if req.Status != currentDevice.Status {
		updates["status"] = req.Status
	}
	if req.DeviceModel != currentDevice.DeviceModel {
		updates["device_model"] = req.DeviceModel
	}
	if req.OperatingSystem != currentDevice.OperatingSystem {
		updates["operating_system"] = req.OperatingSystem
	}

	tx := config.DB.Begin()

	err = tx.Model(&models.PosDevice{}).Where("id = ?", deviceID).Updates(updates).Error

	if err != nil {
		return utils.CapitalizeError(fmt.Sprintf("failed to update device: %v", fmt.Sprintf("error: %v", err)))
	}

	tx.Commit()

	eventservices.RegisterEvent("POS device edited", map[string]interface{}{
		"Pos ID":           deviceID,
		"Serial number":    req.SerialNumber,
		"Business Name":    req.BusinessName,
		"Description":      req.Description,
		"Device Model":     req.DeviceModel,
		"Status":           deviceID,
		"Operating system": req.OperatingSystem,
	})

	return nil
}
