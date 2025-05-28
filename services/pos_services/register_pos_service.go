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

func RegisterPosDevice(req *posdevices.RegisterPosDeviceRequest) error {

	appID, err := uuid.Parse(req.AppId)

	if err != nil {
		utils.Log(slog.LevelError, "Error", "Unable to parse app-Id")
		return utils.CapitalizeError("unable to parse app-Id")
	}

	pos := models.PosDevice{
		ID:                    uuid.New(),
		AppID:                 appID,
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
	}

	result := config.DB.Create(&pos)

	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Unable to register pos device", "Detail", result.Error.Error())
		return utils.CapitalizeError("Unable to register pos device")
	}

	return nil

}
