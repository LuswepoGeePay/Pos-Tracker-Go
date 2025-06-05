package appservices

import (
	"pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	"pos-master/utils"
	"time"
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

func GetAppVersions(req *appPb.GetAppVersionsRequest) (*appPb.GetAppVersionsResponse, error) {

	var app_versions []models.AppVersion

	tx := config.DB.Begin()

	query := tx.Preload("App").Model(&models.AppVersion{})

	var totalApps int64
	err := query.Count(&totalApps).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count app_versions")
	}

	totalPages := int32((totalApps + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&app_versions).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve app_versions")
	}

	pbAppVersions := make([]*appPb.AppVersion, len(app_versions))

	for i, appv := range app_versions {
		pbAppVersions[i] = &appPb.AppVersion{
			AppId:          appv.AppID.String(),
			VersionId:      appv.ID.String(),
			BuildNumber:    appv.BuildNumber,
			ReleaseNotes:   appv.ReleaseNotes,
			VersionNumber:  appv.VersionNumber,
			Checksum:       appv.CheckSum,
			ReleasedAt:     appv.ReleasedAt.Format(time.RFC3339),
			IsActive:       appv.IsActive,
			IsLatestStable: appv.IsLatestStable,
			FilePath:       appv.FilePath,
			FileSizeBytes:  appv.FileSizeBytes,
			AppName:        appv.App.Name,
		}
	}

	return &appPb.GetAppVersionsResponse{
		AppVersion:  pbAppVersions,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}
