package eventservices

import (
	"encoding/json"
	"log/slog"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func RegisterEvent(title string, metaData map[string]interface{}) {

	metaJson, err := json.Marshal(metaData)

	if err != nil {
		utils.Log(slog.LevelError, "error", "unable to marshal metadata", "detail", err.Error())
		return
	}

	event := models.Event{
		ID:    uuid.New(),
		Title: title,

		EventMetaData: datatypes.JSON(metaJson),
	}

	tx := database.DB.Begin()

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "error", "unable to create event", "detail", err.Error())
		return
	}

	tx.Commit()
}
