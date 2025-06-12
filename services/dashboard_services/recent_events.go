package dashboardservices

import (
	"encoding/json"
	"log/slog"
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/dashboard"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func GetRecentEvents() {

}

func GetEvents(req *dashboard.GetEventsRequest) (*dashboard.GetEventsResponse, error) {
	var events []models.Event

	tx := config.DB.Begin()

	query := tx.Model(&models.Event{})

	var totalEvents int64
	err := query.Count(&totalEvents).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count events")
	}

	totalPages := int32((totalEvents + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&events).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve events")
	}

	pbEvents := make([]*dashboard.Event, len(events))

	for i, event := range events {
		pbEvents[i] = &dashboard.Event{
			EventId:  event.ID.String(),
			Title:    event.Title,
			Metadata: string(event.EventMetaData),
			Date:     event.CreatedAt.Format(time.RFC3339),
		}
	}

	return &dashboard.GetEventsResponse{
		Event:       pbEvents,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}

func RegisterEvent(title string, metaData map[string]interface{}) {
	metaJson, err := json.Marshal(metaData)
	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to marshal metadata", "detail", err.Error())
		return
	}

	event := models.Event{
		ID:            uuid.New(),
		Title:         title,
		EventMetaData: datatypes.JSON(metaJson),
	}

	tx := config.DB.Begin()

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "error", "unable to create event", "detail", err.Error())
		return
	}

	tx.Commit()
}
