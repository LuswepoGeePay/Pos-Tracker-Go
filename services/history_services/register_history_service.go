package historyservices

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RegisterNewLocationHistory(req *posdevices.RegisterLocationHistoryRequest) error {

	posID, err := uuid.Parse(req.PosdeviceId)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to parse POS ID", "data", req)
		return utils.CapitalizeError(fmt.Sprintf("failed to parse ID %s", fmt.Sprintf("error: %v", err)))
	}

	var device models.PosDevice

	result := database.DB.Preload("Business").Where("id = ?", posID).Find(&device)
	if result.Error != nil {
		utils.Log(slog.LevelError, "❌error", "Failed to retrieve pos device", "detailed error", result.Error, "request", req)
		return utils.CapitalizeError(fmt.Sprintf("failed to retrieve pos device from ID %s", result.Error))

	}

	locationData, err := GetRegion(req.Longitude, req.Latitude)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve location information", "data", req)
		return utils.CapitalizeError(utils.FormatError("failed to retrieve location information", result.Error))

	}
	newLocation := models.LocationHistory{
		ID:          uuid.New(),
		PosDeviceID: device.ID,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Accuracy:    req.Accuracy,
		TimeStamp:   time.Now(),
		City:        locationData.City,
		RegionName:  locationData.Region,
		IpAddress:   req.IpAddress,
		Entity:      device.Entity,
	}

	// Start transaction for location history
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Log(slog.LevelError, "❌error", "Failed to start transaction", "detail", tx.Error)
		return utils.CapitalizeError(fmt.Sprintf("failed to start transaction: %v", tx.Error))
	}

	result = tx.Create(&newLocation)
	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "Failed to add new location", "location", newLocation)
		return utils.CapitalizeError(fmt.Sprintf("failed to create location: %s", result.Error.Error()))
	}

	floatLongitude, err := strconv.ParseFloat(req.Longitude, 64)
	if err != nil {
		utils.Log(slog.LevelError, "❌error", "Failed to parse longitude", "detail", err, "data", req)
		return utils.CapitalizeError(fmt.Sprintf("failed to parse longitude: %v", err))
	}
	floatLatitude, err := strconv.ParseFloat(req.Latitude, 64)
	if err != nil {
		utils.Log(slog.LevelError, "❌error", "Failed to parse latitude", "detail", err, "data", req)
		return utils.CapitalizeError(fmt.Sprintf("failed to parse latitude: %v", err))
	}

	// Update location on device table within same transaction
	if err := UpdateLocationOnDeviceTableTx(tx, posID.String(), float32(floatLatitude), float32(floatLongitude)); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌error", "Failed to commit transaction", "detail", err)
		return utils.CapitalizeError(fmt.Sprintf("failed to commit transaction: %v", err))
	}

	eventservices.RegisterEvent("A device has pinged", map[string]interface{}{
		"Pos ID":    req.PosdeviceId,
		"longitude": req.Longitude,
		"latitude":  req.Latitude,
		"city":      locationData.City,
		"region":    locationData.Region})

	utils.Log(slog.LevelInfo, "ℹ️Info", "Successfully added new location", "request", req, "location", newLocation)

	return nil
}

type Location struct {
	Region  string
	City    string
	Country string
}

func GetRegion(longitude, latitude string) (*Location, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%s&lon=%s&zoom=10&addressdetails=1", latitude, longitude)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "POSGOAPP/1.0") // Nominatim requires this

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSM returned status: %d", resp.StatusCode)
	}

	var result struct {
		Address struct {
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"address"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &Location{
		City:    result.Address.City,
		Region:  result.Address.State,
		Country: result.Address.Country,
	}, nil
}

func UpdateLocationOnDeviceTable(posID, latitude, longitude string) error {

	updates := map[string]interface{}{}

	if latitude != "" {
		updates["last_known_latitude"] = latitude
	}

	if longitude != "" {
		updates["last_known_longitude"] = longitude
	}

	updates["location_last_updated_at"] = time.Now()

	tx := database.DB.Begin()

	if err := tx.Model(&models.PosDevice{}).Preload("Business").Where("id = ?", posID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError("unable to update pos location")
	}
	tx.Commit()
	return nil

}

// UpdateLocationOnDeviceTableTx updates location within existing transaction
func UpdateLocationOnDeviceTableTx(tx *gorm.DB, posID string, latitude float32, longitude float32) error {
	updates := map[string]interface{}{
		"last_known_latitude":      latitude,
		"last_known_longitude":     longitude,
		"location_last_updated_at": time.Now(),
	}

	if err := tx.Model(&models.PosDevice{}).Where("id = ?", posID).Updates(updates).Error; err != nil {
		return utils.CapitalizeError(fmt.Sprintf("failed to update device location: %v", err))
	}

	return nil
}
