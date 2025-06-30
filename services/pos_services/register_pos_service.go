package posservices

import (
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
)

func RegisterPosDevice(req *posdevices.RegisterPosDeviceRequest) (string, error) {

	posDeviceID := uuid.New()

	var business models.Business

	err := config.DB.Where("email = ?", req.Email).Find(&business)

	if err != nil {
		return "", utils.CapitalizeError(utils.FormatError("unable to find that business", err.Error))
	}

	pos := models.PosDevice{
		ID:                    posDeviceID,
		SerialNumber:          req.SerialNumber,
		Name:                  req.Name,
		Description:           req.Description,
		CurrentAppVersion:     req.CurrentAppVersion,
		LastKnownLatitude:     req.LastKnownLatitude,
		LastKnownLongitude:    req.LastKnownLongitude,
		DeviceModel:           req.DeviceModel,
		OperatingSystem:       req.OperatingSystem,
		Status:                "online",
		LocationLastUpdatedAt: time.Now(),
		Email:                 req.Email,
		Entity:                business.Name,
		FingerPrint:           req.Fingerprint,
		BusinessID:            business.ID,
	}

	result := config.DB.Create(&pos)

	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Unable to register pos device", "Detail", result.Error.Error())
		return "", utils.CapitalizeError("Unable to register pos device")
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
