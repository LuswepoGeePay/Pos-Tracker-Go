package dashboardservices

import (
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/dashboard"
	"pos-master/utils"
	"time"
)

func GetRecentEvents() {

}

func GetEvents(req *dashboard.GetEventsRequest) (*dashboard.GetEventsResponse, error) {
	var events []models.Event

	tx := database.DB.Begin()

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
	err = query.Order("created_at DESC").Limit(int(req.PageSize)).
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
