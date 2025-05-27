package config

import (
	"fmt"
	"log"
	"pos-master/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	var err error

	dsn := "root@tcp(127.0.0.1:3306)/posmaster?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:password@tcp(127.0.0.1:3306)/hof?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Role{},       // Migrate `Role` first as `User` depends on it
		&models.Permission{}, // Other independent tables can be migrated here
		&models.User{},       // Now migrate `User` as `Role` exists
		&models.App{},
		&models.AppVersion{},
		&models.PosDevice{},
		&models.LocationHistory{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database connected successfully")
}
