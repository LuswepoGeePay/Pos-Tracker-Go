package historyservices

import (
	"fmt"
	"log"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
)

func RegisterNewLocationHistory(req *posdevices.RegisterLocationHistoryRequest) error {

	posID, err := uuid.Parse(req.PosdeviceId)

	log.Println(posID)

	if err != nil {
		utils.Log(slog.LevelError, "Error", "Failed to parse POS ID")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse ID %s", fmt.Sprintf("error: %v", err)))
	}

	var device models.PosDevice

	result := config.DB.Where("id = ?", posID).Find(&device)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Failed to retrieve pos device")
		return utils.CapitalizeError(fmt.Sprintf("failed to retrieve pos device from ID %s", result.Error))

	}

	newLocation := models.LocationHistory{
		ID:          uuid.New(),
		PosDeviceID: posID,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Accuracy:    req.Accuracy,
		TimeStamp:   time.Now(),
		City:        req.City,
		RegionName:  req.Region,
		IpAddress:   req.IpAddress,
		Entity:      device.Entity,
	}

	result = config.DB.Create(&newLocation)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Failed to add new location", "detail")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse ID %s", result.Error.Error()))
	}
	return nil
}
