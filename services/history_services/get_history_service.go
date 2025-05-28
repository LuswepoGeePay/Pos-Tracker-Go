package historyservices

import (
	"pos-master/config"
	"pos-master/models"
	posPb "pos-master/proto/posdevices"
	"pos-master/utils"
	"time"
)

func GetLocationHistory(req *posPb.GetLocationHistorysRequest) (*posPb.GetLocationHistorysResponse, error) {

	var pos_location_history []models.LocationHistory

	tx := config.DB.Begin()

	query := tx.Model(&models.PosDevice{})

	var totalPosDevices int64
	err := query.Count(&totalPosDevices).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count pos device location history")
	}

	totalPages := int32((totalPosDevices + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&pos_location_history).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve pos device location history")
	}

	pbLocationHistory := make([]*posPb.LocationHistory, len(pos_location_history))

	for i, history := range pos_location_history {
		pbLocationHistory[i] = &posPb.LocationHistory{
			Id:          history.ID.String(),
			PosdeviceId: history.PosDeviceID.String(),
			Longitude:   history.Longitude,
			Latitude:    history.Latitude,
			Accuracy:    history.Accuracy,
			Timestamp:   history.TimeStamp.Format(time.RFC3339),
		}
	}

	return &posPb.GetLocationHistorysResponse{
		LocationHistory: pbLocationHistory,
		TotalPages:      totalPages,
		CurrentPage:     req.Page,
		HasMore:         req.Page < totalPages,
	}, nil
}

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
			AppId:               history.AppID.String(),
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
		Posdevices:  pbPosDevices,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}
