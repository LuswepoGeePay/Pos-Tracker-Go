package posservices

import (
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	posPb "pos-master/proto/posdevices"
	"pos-master/utils"
	"time"
)

func GetPosDevices(req *posPb.GetPosDevicesRequest) (*posPb.GetPosDevicesResponse, error) {

	var pos_devices []models.PosDevice

	tx := config.DB.Begin()

	query := tx.Model(&models.PosDevice{})

	var totalPosDevices int64
	err := query.Count(&totalPosDevices).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count pos device pos devices")
	}

	totalPages := int32((totalPosDevices + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&pos_devices).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve pos device pos devices")
	}

	pbPosDevices := make([]*posPb.PosDevice, len(pos_devices))

	for i, history := range pos_devices {
		pbPosDevices[i] = &posPb.PosDevice{
			Id:                  history.ID.String(),
			SerialNumber:        history.SerialNumber,
			Name:                history.Name,
			CurrentAppVersion:   history.CurrentAppVersion,
			LastKnownLatitude:   history.LastKnownLatitude,
			LastKnownLongitude:  history.LastKnownLongitutude,
			Status:              history.Status,
			DeviceModel:         history.DeviceModel,
			OperatingSystem:     history.OperatingSystem,
			Description:         history.Description,
			LocationLastUpdated: history.LocationLastUpdatedAt.Format(time.RFC3339),
		}
	}

	return &posPb.GetPosDevicesResponse{
		Posdevice:   pbPosDevices,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}

func GetPosById(posID string) (*posPb.PosDevice, error) {

	var posdevice models.PosDevice

	tx := config.DB.Begin()

	result := tx.First(&posdevice, "id = ?", posID)

	if result.Error != nil {
		utils.Log(slog.LevelError, "❌ error", result.Error)
		return nil, result.Error
	}

	return &posPb.PosDevice{
		Id:                 posdevice.ID.String(),
		SerialNumber:       posdevice.SerialNumber,
		Name:               posdevice.Name,
		Description:        posdevice.Description,
		CurrentAppVersion:  posdevice.CurrentAppVersion,
		LastKnownLatitude:  posdevice.LastKnownLatitude,
		LastKnownLongitude: posdevice.LastKnownLongitutude,
		Status:             posdevice.Status,
		DeviceModel:        posdevice.DeviceModel,
		OperatingSystem:    posdevice.OperatingSystem,
	}, nil

}
