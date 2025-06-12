package appservices

import (
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	"pos-master/utils"
)

func CheckAppUpdate(req *posdevices.CheckUpdateRequest) (*posdevices.CheckUpdateResponse, error) {

	var latestVersion models.AppVersion

	err := config.DB.Where("is_active = ? AND is_latest_stable = ?", true, true).Order("released_at desc").First(&latestVersion).Error

	if err != nil {
		return nil, utils.CapitalizeError("could not fetch latest app version")
	}

	// if latest.VersionNumber == req.AppVersion && latest.BuildNumber == req.BuildVersion {
	// 	// Already up to date
	// 	return &posdevices.CheckUpdateResponse{
	// 		UpdateAvailable: false,
	// 		LatestVersion:   latest.VersionNumber,
	// 		ReleaseNotes:    latest.ReleaseNotes,
	// 		DownloadUrl:     "",
	// 		Code:            0, // no update
	// 	}, nil
	// }

	if latestVersion.VersionNumber == req.AppVersion {
		return &posdevices.CheckUpdateResponse{
			UpdateAvailable: false,
			LatestVersion:   latestVersion.VersionNumber,
			ReleaseNotes:    latestVersion.ReleaseNotes,
			DownloadUrl:     "",
			Code:            0,
		}, nil
	}

	return &posdevices.CheckUpdateResponse{
		UpdateAvailable: true,
		LatestVersion:   latestVersion.VersionNumber,
		ReleaseNotes:    latestVersion.ReleaseNotes,
		DownloadUrl:     latestVersion.FilePath,
		Code:            0,
	}, nil
}
