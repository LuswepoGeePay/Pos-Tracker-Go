package appservices

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	"pos-master/services/pocketbase"
	"pos-master/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterApp(req *appPb.RegisterAppRequest) error {

	app := models.App{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	result := config.DB.Create(&app)
	if result.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to register app %s", result.Error.Error()))
	}

	return nil

}

func RegisterAppVersion(c *gin.Context, req *appPb.RegisterAppVersionRequest) error {

	token, err := pocketbase.HandlePocketBaseAuth(c)

	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to get pocketbase token", "detail", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to get pocketbase token")
	}

	fileURL, err := pocketbase.HandleUpload(c, token, "file")
	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to upload file to pocketbase", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to upload file to server")
	}

	// file, err := c.FormFile("file")

	// if err != nil {
	// 	return utils.CapitalizeError("no file uploaded")
	// }

	// openedFile, err := file.Open()
	// if err != nil {
	// 	utils.Log(slog.LevelError, "error", "unable to open APK for checksum", fmt.Sprintf("error: %v", err))
	// 	return utils.CapitalizeError("unable to read uploaded APK")
	// }

	// defer openedFile.Close()
	res, err := http.Get(fileURL)
	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to fetch APK from URL", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to fetch uploaded APK")
	}
	defer res.Body.Close()

	hash := sha256.New()
	// if _, err := io.Copy(hash, openedFile); err != nil {
	// 	utils.Log(slog.LevelError, "❌error", "unable to hash APK", fmt.Sprintf("error: %v", err))
	// 	return utils.CapitalizeError("unable to process apk")
	// }

	appID, err := uuid.Parse(req.AppId)

	if err != nil {
		utils.Log(slog.LevelError, "error", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to parse App ID")
	}

	checkSum := hex.EncodeToString(hash.Sum(nil))

	appVersion := models.AppVersion{
		ID:             uuid.New(),
		AppID:          appID,
		ReleaseNotes:   req.ReleaseNotes,
		FilePath:       fileURL,
		FileSizeMBytes: req.FileSizeBytes,
		CheckSum:       checkSum,
		IsActive:       false,
		IsLatestStable: req.IsLatestStable,
		ReleasedAt:     time.Now(),
		VersionNumber:  req.VersionNumber,
	}

	result := config.DB.Create(&appVersion)

	if result.Error != nil {

	}

	return nil

}
