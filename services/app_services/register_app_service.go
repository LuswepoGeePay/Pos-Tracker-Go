package appservices

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	database "pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	eventservices "pos-master/services/event_services"
	"pos-master/services/pocketbase"
	"pos-master/utils"
	"time"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterApp(req *appPb.RegisterAppRequest) error {

	app := models.App{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	// Start transaction for app registration
	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to start transaction: %v", tx.Error))
	}

	result := tx.Create(&app)
	if result.Error != nil {
		tx.Rollback()
		return utils.CapitalizeError(fmt.Sprintf("Unable to register app %s", result.Error.Error()))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	eventservices.RegisterEvent("New App Registered", map[string]interface{}{
		"App Name":    req.Name,
		"Description": req.Description,
	})

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
	// Read the body once and hash it. Note: since we only need the hash, we can pipe it.
	// But we also need to close res.Body.
	if _, err := io.Copy(hash, res.Body); err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to hash APK", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to process apk")
	}

	appID, err := uuid.Parse(req.AppId)

	if err != nil {
		utils.Log(slog.LevelError, "error", fmt.Sprintf("error: %v", err))
		return utils.CapitalizeError("unable to parse App ID")
	}

	var terminalTypeID *uuid.UUID
	if req.TerminalTypeId != "" {
		parsedID, err := uuid.Parse(req.TerminalTypeId)
		if err == nil {
			terminalTypeID = &parsedID
		}
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
		TerminalTypeID: terminalTypeID,
	}

	// Start transaction for app version registration
	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to start transaction: %v", tx.Error))
	}

	result := tx.Create(&appVersion)
	if result.Error != nil {
		tx.Rollback()
		return utils.CapitalizeError(utils.FormatError("unable to create app version", result.Error))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	eventservices.RegisterEvent("New App Version Registered", map[string]interface{}{
		"App Id":         appID,
		"Release notes":  req.ReleaseNotes,
		"Version number": req.VersionNumber,
		"File path":      fileURL,
	})

	return nil

}
