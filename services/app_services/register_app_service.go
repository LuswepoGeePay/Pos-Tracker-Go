package appservices

import (
	"fmt"
	"pos-master/config"
	"pos-master/models"
	appPb "pos-master/proto/app"
	"pos-master/utils"

	"github.com/google/uuid"
)

func RegisterApp(req *appPb.RegisterAppRequest) error {

	app := models.App{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	result := config.DB.Create(&app)
	if result.Error != nil {
		return utils.CapitalizeError(fmt.Sprintf("Unable to register app %s", result.Error.Error()))
	}

	return nil

}
