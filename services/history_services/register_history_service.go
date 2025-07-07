package historyservices

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/posdevices"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
)

func RegisterNewLocationHistory(req *posdevices.RegisterLocationHistoryRequest) error {

	posID, err := uuid.Parse(req.PosdeviceId)

	log.Println(posID)

	if err != nil {
		utils.Log(slog.LevelError, "Error", "Failed to parse POS ID")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse ID %s", fmt.Sprintf("error: %v", err)))
	}

	var device models.PosDevice

	result := config.DB.Where("id = ?", posID).Find(&device)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Failed to retrieve pos device")
		return utils.CapitalizeError(fmt.Sprintf("failed to retrieve pos device from ID %s", result.Error))

	}

	locationData, err := GetRegion(req.Longitude, req.Latitude)

	if err != nil {
		utils.Log(slog.LevelError, "Error", "unable to retrieve location information")
		return utils.CapitalizeError(utils.FormatError("failed to retrieve location information", result.Error))

	}
	newLocation := models.LocationHistory{
		ID:          uuid.New(),
		PosDeviceID: posID,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Accuracy:    req.Accuracy,
		TimeStamp:   time.Now(),
		City:        locationData.City,
		RegionName:  locationData.Region,
		IpAddress:   req.IpAddress,
		Entity:      device.Entity,
	}

	result = config.DB.Create(&newLocation)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Failed to add new location", "detail")
		return utils.CapitalizeError(fmt.Sprintf("failed to parse ID %s", result.Error.Error()))
	}

	_ = UpdateLocationOnDeviceTable(posID.String(), req.Latitude, req.Longitude)

	eventservices.RegisterEvent("A device has pinged", map[string]interface{}{
		"Pos ID":    req.PosdeviceId,
		"longitude": req.Longitude,
		"latitude":  req.Latitude,
		"city":      locationData.City,
		"region":    locationData.Region})
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

	tx := config.DB.Begin()

	if err := tx.Model(&models.PosDevice{}).Where("id = ?", posID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return utils.CapitalizeError("unable to update pos location")
	}
	tx.Commit()
	return nil

}
