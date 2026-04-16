package posservices

import (
	"fmt"
	"log/slog"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func RegisterHeartBeat(req *models.HeartBeatRequest) error {
	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse pos device id", "data", req)
		return utils.CapitalizeError("invalid device ID")
	}

	// Try to parse the timestamp with multiple common layouts
	var ts time.Time
	layouts := []string{
		time.RFC3339,               // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05.999999", // 2026-04-16T13:43:32.301964
		"2006-01-02 15:04:05",
	}

	parsed := false
	for _, layout := range layouts {
		if t, err := time.Parse(layout, req.Timestamp); err == nil {
			ts = t
			parsed = true
			break
		}
	}

	if !parsed {
		// If we couldn't parse it, use current time as fallback
		ts = time.Now()
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var device models.PosDevice
	result := tx.Where("id = ?", deviceID).First(&device)
	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to find pos device with this id", "data", req)
		return utils.CapitalizeError("cannot find device with this ID")
	}

	updateResult := tx.Model(&device).Updates(map[string]interface{}{
		"location_last_updated_at": ts,
		"status":                   "online",
	})

	if updateResult.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "unable to update pos device status", "data", req, "detail", updateResult.Error)

		return utils.CapitalizeError("unable to update device status")
	}

	tx.Commit()

	utils.Log(slog.LevelInfo, "ℹ️info", "successfully registered heart beat", "data", req)

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

	result := database.DB.Exec(query, threshold)
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

	result := database.DB.Exec(query, threshold)
	fmt.Printf("[%s] Fallback marked %d devices as online from LocationHistory\n", time.Now().Format(time.RFC3339), result.RowsAffected)
}
func StartCronJobs() {
	c := cron.New()
	// Runs every 30 minutes. Use cron expression format.
	_, err := c.AddFunc("*/30 * * *  *", MarkOnlineFromLocationHistory)
	if err != nil {
		panic(err)
	}
	_, err = c.AddFunc("*/30 * * * *", MarkOfflineDevices)

	if err != nil {
		panic(err)
	}
	c.Start()
}
