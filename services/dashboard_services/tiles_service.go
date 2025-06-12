package dashboardservices

import (
	"errors"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/dashboard"
	"pos-master/utils"
	"time"

	"gorm.io/gorm"
)

func GetTileInfo() (*dashboard.Tile, error) {

	var totalPosDevices int64
	var totalActiveDevices int64
	var totalOfflineDevices int64
	var totalRegisteredApps int64
	var latestAppVersion string
	var totalLocationPings int64

	tx := config.DB

	//pos devices
	devicesQuery := tx.Model(&models.PosDevice{})

	err := devicesQuery.Count(&totalPosDevices).Error
	if err != nil {
		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count pos devices", err))
		return nil, utils.CapitalizeError(utils.FormatError("unable to count pos devices", err))
	}

	//active devices
	activeDevicesQuery := tx.Model(&models.PosDevice{})

	err = activeDevicesQuery.Where("status = ?", "online").Count(&totalActiveDevices).Error
	if err != nil {
		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count active pos devices", err))
		return nil, utils.CapitalizeError(utils.FormatError("unable to count active pos devices", err))
	}

	//offline devices
	offlineDevicesQuery := tx.Model(&models.PosDevice{})

	err = offlineDevicesQuery.Where("status = ? or status = ?", "offline", "inactive").Count(&totalOfflineDevices).Error
	if err != nil {
		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count offline pos devices", err))
		return nil, utils.CapitalizeError(utils.FormatError("unable to count offline pos devices", err))
	}

	//apps
	appsQuery := tx.Model(&models.App{})

	err = appsQuery.Count(&totalRegisteredApps).Error
	if err != nil {
		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count registered apps", err))
		return nil, utils.CapitalizeError(utils.FormatError("unable to count registered apps", err))
	}

	//locations

	now := time.Now()

	location, _ := time.LoadLocation("Africa/Lusaka") // adjust based on deployment

	// Start of Today at 01:00 AM
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, location)
	endOfToday := startOfToday.AddDate(0, 0, 1)

	locationsQuery := tx.Model(&models.LocationHistory{}).
		Where("created_at >= ? AND created_at < ?", startOfToday, endOfToday)

	err = locationsQuery.Count(&totalLocationPings).Error
	if err != nil {
		utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to count location pings today", err))
		return nil, utils.CapitalizeError(utils.FormatError("unable to count location pings today", err))
	}

	//app versions
	var version models.AppVersion

	err = tx.Where("is_active = ? AND is_latest_stable = ?", true, true).
		Order("released_at desc").
		First(&version).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// This is fine — we’ll just return empty version string
			latestAppVersion = "N/A"
		} else {
			utils.Log(slog.LevelError, "❌ error", utils.FormatError("unable to fetch latest app version", err))
			return nil, utils.CapitalizeError(utils.FormatError("unable to fetch latest app version", err))
		}
	} else {
		latestAppVersion = version.VersionNumber
	}

	return &dashboard.Tile{
		PosDevices:       int32(totalPosDevices),
		ActiveDevices:    int32(totalActiveDevices),
		OfflineDevices:   int32(totalOfflineDevices),
		Apps:             int32(totalRegisteredApps),
		AppVersion:       latestAppVersion,
		LocationsTracked: int32(totalLocationPings),
	}, nil
}
