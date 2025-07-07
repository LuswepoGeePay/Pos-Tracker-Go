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
	threshold := time.Now().Add(-168 * time.Hour)

	query := `
		UPDATE pos_devices pd
		JOIN (
			SELECT lh.pos_device_id, MAX(lh.time_stamp) AS latest_ts
			FROM location_histories lh
			GROUP BY lh.pos_device_id
			HAVING latest_ts < ?
		) AS old ON pd.id = old.pos_device_id
		SET pd.status = 'offline'
		WHERE pd.status != 'offline'
	`

	result := config.DB.Exec(query, threshold)
	if result.Error != nil {
		fmt.Printf("Error during fallback offline update: %v\n", result.Error)
	}
	fmt.Printf("[%s] Fallback marked %d devices as offline from LocationHistory\n",
		time.Now().Format(time.RFC3339), result.RowsAffected)
}

func MarkOnlineFromLocationHistory() {
	threshold := time.Now().Add(-24 * time.Hour)

	// Raw SQL can be more efficient here depending on the size of the history table
	query := `
		UPDATE pos_devices pd
		JOIN (
			SELECT lh.pos_device_id, MAX(lh.time_stamp) AS latest_ts
			FROM location_histories lh
			GROUP BY lh.pos_device_id
			HAVING latest_ts >= ?
		) AS recent ON pd.id = recent.pos_device_id
		SET pd.status = 'online'
		WHERE pd.status != 'online'
	`

	result := config.DB.Exec(query, threshold)
	fmt.Printf("[%s] Fallback marked %d devices as online from LocationHistory\n", time.Now().Format(time.RFC3339), result.RowsAffected)
}
func StartCronJobs() {
	c := cron.New()
	// Runs every 10 minutes. Use cron expression format.
	_, err := c.AddFunc("*/10 * * *  *", MarkOnlineFromLocationHistory)
	if err != nil {
		panic(err)
	}
	_, err = c.AddFunc("*/10 * * * *", MarkOfflineDevices)

	if err != nil {
		panic(err)
	}
	c.Start()
}
