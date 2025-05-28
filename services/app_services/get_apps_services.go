package appservices

import (
	"pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	"pos-master/utils"
)

func GetApps(req *appPb.GetAppsRequest) (*appPb.GetAppsResponse, error) {

	var apps []models.App

	tx := config.DB.Begin()

	query := tx.Model(&models.App{})

	var totalApps int64
	err := query.Count(&totalApps).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count apps")
	}

	totalPages := int32((totalApps + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&apps).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve apps")
	}

	pbApps := make([]*appPb.App, len(apps))

	for i, app := range apps {
		pbApps[i] = &appPb.App{
			Id:          app.ID.String(),
			Name:        app.Name,
			Description: app.Description,
		}
	}

	return &appPb.GetAppsResponse{
		App:         pbApps,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}
