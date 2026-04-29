package dashboardservices

import (
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/dashboard"
	"time"
)

// func GetTileInfo() (*dashboard.Tile, error) {

// 	var totalPosDevices int64
// 	var totalActiveDevices int64
// 	var totalOfflineDevices int64
// 	var totalRegisteredApps int64
// 	var latestAppVersion string
// 	var totalLocationPings int64

// 	tx := database.DB

// 	//pos devices
// 	devicesQuery := tx.Model(&models.PosDevice{})

// 	err := devicesQuery.Count(&totalPosDevices).Error
// 	if err != nil {
// 		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count pos devices", err))
// 		return nil, utils.CapitalizeError(utils.FormatError("unable to count pos devices", err))
// 	}

// 	//active devices
// 	activeDevicesQuery := tx.Model(&models.PosDevice{})

// 	err = activeDevicesQuery.Where("status = ?", "online").Count(&totalActiveDevices).Error
// 	if err != nil {
// 		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count active pos devices", err))
// 		return nil, utils.CapitalizeError(utils.FormatError("unable to count active pos devices", err))
// 	}

// 	//offline devices
// 	offlineDevicesQuery := tx.Model(&models.PosDevice{})

// 	err = offlineDevicesQuery.Where("status = ? or status = ?", "offline", "inactive").Count(&totalOfflineDevices).Error
// 	if err != nil {
// 		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count offline pos devices", err))
// 		return nil, utils.CapitalizeError(utils.FormatError("unable to count offline pos devices", err))
// 	}

// 	//apps
// 	appsQuery := tx.Model(&models.App{})

// 	err = appsQuery.Count(&totalRegisteredApps).Error
// 	if err != nil {
// 		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count registered apps", err))
// 		return nil, utils.CapitalizeError(utils.FormatError("unable to count registered apps", err))
// 	}

// 	//locations

// 	now := time.Now()

// 	location, _ := time.LoadLocation("Africa/Lusaka") // adjust based on deployment

// 	// Start of Today at 01:00 AM
// 	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, location)
// 	endOfToday := startOfToday.AddDate(0, 0, 1)

// 	locationsQuery := tx.Model(&models.LocationHistory{}).
// 		Where("created_at >= ? AND created_at < ?", startOfToday, endOfToday)

// 	err = locationsQuery.Count(&totalLocationPings).Error
// 	if err != nil {
// 		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count location pings today", err))
// 		return nil, utils.CapitalizeError(utils.FormatError("unable to count location pings today", err))
// 	}

// 	//app versions
// 	var version models.AppVersion

// 	err = tx.Where("is_active = ? AND is_latest_stable = ?", true, true).
// 		Order("released_at desc").
// 		First(&version).Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// This is fine — we’ll just return empty version string
// 			latestAppVersion = "N/A"
// 		} else {
// 			utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to fetch latest app version", err))
// 			return nil, utils.CapitalizeError(utils.FormatError("unable to fetch latest app version", err))
// 		}
// 	} else {
// 		latestAppVersion = version.VersionNumber
// 	}

// 	return &dashboard.Tile{
// 		PosDevices:       int32(totalPosDevices),
// 		ActiveDevices:    int32(totalActiveDevices),
// 		OfflineDevices:   int32(totalOfflineDevices),
// 		Apps:             int32(totalRegisteredApps),
// 		AppVersion:       latestAppVersion,
// 		LocationsTracked: int32(totalLocationPings),
// 	}, nil
// }

func GetTileInfo() (*dashboard.Tile, error) {
	tx := database.DB

	var result struct {
		Total   int64
		Active  int64
		Offline int64
	}

	err := tx.Raw(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN status='online' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status IN ('offline','inactive') THEN 1 ELSE 0 END) as offline
		FROM pos_devices
		WHERE deleted_at IS NULL
	`).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	var apps int64
	tx.Model(&models.App{}).Count(&apps)

	location, err := time.LoadLocation("Africa/Lusaka")
	if err != nil {
		location = time.UTC
	}

	now := time.Now().In(location)

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	end := start.Add(24 * time.Hour)

	var locations int64
	tx.Model(&models.LocationHistory{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&locations)

	var version models.AppVersion
	latest := "N/A"

	err = tx.Where("is_active = ? AND is_latest_stable = ?", true, true).
		Order("released_at desc").
		First(&version).Error

	if err == nil {
		latest = version.VersionNumber
	}

	return &dashboard.Tile{
		PosDevices:       int32(result.Total),
		ActiveDevices:    int32(result.Active),
		OfflineDevices:   int32(result.Offline),
		Apps:             int32(apps),
		AppVersion:       latest,
		LocationsTracked: int32(locations),
	}, nil
}
