package seeders

import (
	"fmt"
	"pos-master/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedTerminalTypes(db *gorm.DB) error {
	deviceTypes := []models.TerminalType{
		{
			ID:            uuid.New(),
			Name:          "Trendit Terminal",
			TerminalModel: "S680",
		},
		{
			ID:            uuid.New(),
			Name:          "NewLand Terminal",
			TerminalModel: "N950",
		},
	}

	for _, deviceType := range deviceTypes {
		var existing models.TerminalType

		err := db.Where("name = ?", deviceType.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&deviceType).Error; err != nil {
				return err
			}
			fmt.Println("Seeded device type:", deviceType.Name)
		}
	}

	return nil
}
