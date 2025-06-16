package appservices

import (
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
)

func CheckAppUpdate(req *posdevices.CheckUpdateRequest) (*posdevices.CheckUpdateResponse, error) {

	var latestVersion models.AppVersion

	err := config.DB.Where("is_active = ? AND is_latest_stable = ?", true, true).Order("released_at desc").First(&latestVersion).Error

	if err != nil {
		return nil, utils.CapitalizeError("could not fetch latest app version")
	}

	var posDevice models.PosDevice

	tx := config.DB
	err = tx.Where("id = ?", req.PosdeviceId).Find(&posDevice).Error

	if err != nil {
		return nil, utils.CapitalizeError(utils.FormatError("unable to find pos device with this ID", err))
	}
	if posDevice.CurrentAppVersion != latestVersion.VersionNumber {
		result := tx.Model(&models.PosDevice{}).Where("id = ?", req.PosdeviceId).Update("current_app_version", latestVersion.VersionNumber)
		if result.Error != nil {
			return nil, utils.CapitalizeError(utils.FormatError("unable to update pos device app version", result.Error))
		}
	}

	if latestVersion.VersionNumber == req.AppVersion {
		return &posdevices.CheckUpdateResponse{
			UpdateAvailable: false,
			LatestVersion:   latestVersion.VersionNumber,
			ReleaseNotes:    latestVersion.ReleaseNotes,
			DownloadUrl:     "",
			Code:            0,
		}, nil
	}

	eventservices.RegisterEvent("Pos Device checked for an update", map[string]interface{}{
		"pos_device_id": req.PosdeviceId,
		"version":       req.AppVersion,
	})

	return &posdevices.CheckUpdateResponse{
		UpdateAvailable: true,
		LatestVersion:   latestVersion.VersionNumber,
		ReleaseNotes:    latestVersion.ReleaseNotes,
		DownloadUrl:     latestVersion.FilePath,
		Code:            0,
	}, nil
}
