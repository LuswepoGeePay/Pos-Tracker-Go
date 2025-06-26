package posservices

import (
	"fmt"
	"pos-master/config"
	"pos-master/models"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func RegisterHeartBeat(req *models.HeartBeatRequest) error {
	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		return utils.CapitalizeError("invalid device ID")
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var device models.PosDevice
	result := tx.Where("id = ?", deviceID).First(&device)
	if result.Error != nil {
		tx.Rollback()
		return utils.CapitalizeError("cannot find device with this ID")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return utils.CapitalizeError("device not found")
	}

	updateResult := tx.Model(&device).Updates(map[string]interface{}{
		"location_last_updated_at": req.Timestamp,
		"status":                   "online",
	})

	if updateResult.Error != nil {
		tx.Rollback()
		return utils.CapitalizeError("unable to update device status")
	}

	tx.Commit()
	return nil
}

func MarkOfflineDevices() {
	threshold := time.Now().Add(-3 * time.Hour)

	result := config.DB.Model(&models.PosDevice{}).Where("location_last_updated_at < ?", threshold).Update("status", "offline")
	fmt.Printf("[%s] Marked %d devices as offline\n", time.Now().Format(time.RFC3339), result.RowsAffected)

}

func StartCronJobs() {
	c := cron.New()
	// Runs every 10 minutes. Use cron expression format.
	_, err := c.AddFunc("*/10 * * * *", MarkOfflineDevices)
	if err != nil {
		panic(err)
	}

	c.Start()
}
