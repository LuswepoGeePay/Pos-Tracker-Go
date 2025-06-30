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

	query := tx.Preload("PosDevice").Model(&models.LocationHistory{})

	var totalLocations int64
	err := query.Count(&totalLocations).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count pos device location history")
	}

	totalPages := int32((totalLocations + int64(req.PageSize) - 1) / int64(req.PageSize))
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
			Id:           history.ID.String(),
			PosdeviceId:  history.PosDeviceID.String(),
			Longitude:    history.Longitude,
			Latitude:     history.Latitude,
			Accuracy:     history.Accuracy,
			Timestamp:    history.TimeStamp.Format(time.RFC3339),
			DeviceName:   history.PosDevice.Name,
			BusinessName: history.Entity,
			Region:       history.RegionName,
		}
	}

	return &posPb.GetLocationHistorysResponse{
		LocationHistory: pbLocationHistory,
		TotalPages:      totalPages,
		CurrentPage:     req.Page,
		HasMore:         req.Page < totalPages,
	}, nil
}
