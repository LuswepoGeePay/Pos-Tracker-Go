package businessservices

import (
	"pos-master/config"
	"pos-master/models"
	"pos-master/proto/business"
	"pos-master/utils"
	"time"

	"github.com/google/uuid"
)

func GetBusinesses(req *business.GetBusinessesRequest) (*business.GetBusinessesResponse, error) {

	var businesses []models.Business

	tx := config.DB.Begin()

	query := tx.Model(&models.Business{})

	var totalBusinesss int64
	err := query.Count(&totalBusinesss).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count businesses")
	}

	totalPages := int32((totalBusinesss + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&businesses).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve businesses")
	}

	pbBusinesss := make([]*business.Business, len(businesses))

	for i, bs := range businesses {
		pbBusinesss[i] = &business.Business{
			Id:      bs.ID.String(),
			Name:    bs.Name,
			Status:  bs.Status,
			Email:   bs.Email,
			Phone:   bs.Phone,
			Address: bs.Address,
		}
	}

	return &business.GetBusinessesResponse{
		Business:    pbBusinesss,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
		Count:       int32(totalBusinesss),
	}, nil

}

func GetBusinessById(businessID string) (*business.Business, error) {

	parsedID, err := uuid.Parse(businessID)

	if err != nil {
		return nil, utils.CapitalizeError("inavlid business ID")
	}

	var bmodel models.Business

	err = config.DB.Where("id = ?", parsedID).Find(&bmodel).Error

	if err != nil {
		return nil, utils.CapitalizeError(utils.FormatError("unable to find business", err))
	}

	// Fetch POS devices for this business
	var deviceModels []models.PosDevice
	if err := config.DB.Preload("Business").Where("business_id = ?", parsedID).Find(&deviceModels).Error; err != nil {
		return nil, utils.CapitalizeError(utils.FormatError("unable to fetch pos devices", err))
	}

	pbDevices := make([]*business.PosDevice, 0, len(deviceModels))

	for _, d := range deviceModels {
		pbDevices = append(pbDevices, &business.PosDevice{
			Id:                  d.ID.String(),
			Name:                d.Name,
			Description:         d.Description,
			BusinessName:        d.Business.Name,
			CurrentAppVersion:   d.CurrentAppVersion,
			SerialNumber:        d.SerialNumber,
			Status:              d.Status,
			OperatingSystem:     d.OperatingSystem,
			LastKnownLatitude:   d.LastKnownLatitude,
			LastKnownLongitude:  d.LastKnownLongitude,
			LocationLastUpdated: d.LocationLastUpdatedAt.Format(time.RFC3339),
			DeviceModel:         d.DeviceModel,
			Fingerprint:         d.FingerPrint,
		})
	}
	return &business.Business{
		Id:           bmodel.ID.String(),
		BusinessLogo: bmodel.BusinessLogo,
		Name:         bmodel.Name,
		Address:      bmodel.Address,
		Status:       bmodel.Status,
		Device:       pbDevices,
	}, nil
}
