package posservices

import (
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
)

func RegisterPosDevice(req *posdevices.RegisterPosDeviceRequest) (string, error) {

	posDeviceID := uuid.New()

	pos := models.PosDevice{
		ID:                    posDeviceID,
		SerialNumber:          req.SerialNumber,
		Name:                  req.Name,
		Description:           req.Description,
		CurrentAppVersion:     req.CurrentAppVersion,
		LastKnownLatitude:     req.LastKnownLatitude,
		LastKnownLongitutude:  req.LastKnownLongitude,
		DeviceModel:           req.DeviceModel,
		OperatingSystem:       req.OperatingSystem,
		Status:                "online",
		LocationLastUpdatedAt: time.Now(),
		Email:                 req.Email,
	}

	result := config.DB.Create(&pos)

	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Unable to register pos device", "Detail", result.Error.Error())
		return "", utils.CapitalizeError("Unable to register pos device")
	}

	return posDeviceID.String(), nil

}
