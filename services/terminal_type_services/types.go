package terminaltypeservices

import (
	"fmt"
	"log/slog"
	database "pos-master/config"
	"pos-master/models"
	"pos-master/proto/terminaltype"
	"pos-master/utils"

	"github.com/google/uuid"
)

func CreateTerminalType(req *terminaltype.RegisterTerminalTypeRequest) error {

	tType := models.TerminalType{
		ID:            uuid.New(),
		Name:          req.Name,
		TerminalModel: req.TerminalModel,
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Create(&tType)
	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error", "Unable to create terminal type", "Detail", result.Error.Error(), "data", req)
		return utils.CapitalizeError("Unable to create terminal type")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error", "Failed to commit transaction", "Detail", err.Error())
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	return nil
}

func GetTerminalTypes(req *terminaltype.GetTerminalTypesRequest) (*terminaltype.GetTerminalTypesResponse, error) {

	var terminalTypes []models.TerminalType
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := tx.Model(&models.TerminalType{})

	var totalTerminalTypes int64
	err := query.Count(&totalTerminalTypes).Error
	if err != nil {
		return nil, utils.CapitalizeError("failed to count terminal types")
	}

	totalPages := int32((totalTerminalTypes + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	if req.SearchQuery != "" {
		query = query.Where("name LIKE ? OR terminal_model LIKE ?", "%"+req.SearchQuery+"%", "%"+req.SearchQuery+"%")
	}

	result := query.Order("created_at DESC").Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&terminalTypes)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Unable to get terminal types", "Detail", result.Error.Error())
		return nil, utils.CapitalizeError("Unable to get terminal types")
	}

	pbTerminalTypes := make([]*terminaltype.TerminalType, len(terminalTypes))
	for i, terminalType := range terminalTypes {
		pbTerminalTypes[i] = &terminaltype.TerminalType{
			Id:            terminalType.ID.String(),
			Name:          terminalType.Name,
			TerminalModel: terminalType.TerminalModel,
		}
	}
	return &terminaltype.GetTerminalTypesResponse{
		TerminalTypes: pbTerminalTypes,
		TotalPages:    totalPages,
		CurrentPage:   req.Page,
		HasMore:       req.Page < totalPages,
		Count:         int32(totalTerminalTypes),
	}, nil
}

func EditTerminalType(req *terminaltype.EditTerminalTypeRequest) error {

	upddates := make(map[string]interface{})

	if req.Name != "" {
		upddates["name"] = req.Name
	}

	if req.TerminalModel != "" {
		upddates["terminal_model"] = req.TerminalModel
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("id = ?", req.Id).Updates(&upddates)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Unable to update terminal type", "Detail", result.Error.Error(), "data", req)
		return utils.CapitalizeError("Unable to update terminal type")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error", "Failed to commit transaction", "Detail", err.Error())
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	return nil
}

func DeleteTerminalType(typeId string) error {

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("id = ?", typeId).Delete(&models.TerminalType{})
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Unable to delete terminal type", "Detail", result.Error.Error(), "data", typeId)
		return utils.CapitalizeError("Unable to delete terminal type")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error", "Failed to commit transaction", "Detail", err.Error())
		return utils.CapitalizeError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	return nil
}
