package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	FullName string    `gorm:"default:null"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	RoleID   uuid.UUID `gorm:"not null"`
	Role     Role      `gorm:"foreignKey:RoleID"`
	Status   bool      `gorm:"default:false"`
	gorm.Model
	// UserID          uuid.UUID `gorm:"type:uuid;not null;unique"`
	// User            User      `gorm:"foreignKey:UserID"`
}

type App struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"default:null"`
	Description string    `gorm:"default:null"`
	gorm.Model
}

type AppVersion struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	AppID          uuid.UUID `gorm:"not null"`
	App            App       `gorm:"foreignKey:AppID"`
	VersionNumber  string    `gorm:"default:null"`
	ReleaseNotes   string    `gorm:"default:null"`
	FilePath       string    `gorm:"default:null"`
	FileSizeMBytes string    `gorm:"default:null"`
	CheckSum       string    `gorm:"default:null"`
	IsActive       bool      `gorm:"default:false"`
	IsLatestStable bool      `gorm:"default:false"`
	ReleasedAt     time.Time
	gorm.Model
}

type PosDevice struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key"`
	SerialNumber          string    `gorm:"default:null"`
	Entity                string    `gorm:"default:null"`
	Name                  string    `gorm:"default:null"`
	Description           string    `gorm:"default:null"`
	CurrentAppVersion     string    `gorm:"default:null"`
	LastKnownLatitude     string    `gorm:"default:null"`
	LastKnownLongitutude  string    `gorm:"default:null"`
	Status                string    `gorm:"default:null"`
	DeviceModel           string    `gorm:"default:null"`
	OperatingSystem       string    `gorm:"default:null"`
	Email                 string    `gorm:"default:null"`
	LocationLastUpdatedAt time.Time
	gorm.Model
}

type LocationHistory struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	PosDeviceID uuid.UUID `gorm:"not null"`
	PosDevice   PosDevice `gorm:"foreignKey:PosDeviceID"`
	Longitude   string    `gorm:"default:null"`
	Latitude    string    `gorm:"default:null"`
	Accuracy    string    `gorm:"default:null"`
	City        string    `gorm:"default:null"`
	IpAddress   string    `gorm:"default:null"`
	RegionName  string    `gorm:"default:null"`
	Entity      string    `gorm:"default:null"`
	TimeStamp   time.Time
	gorm.Model
}

type Permission struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	Name string    `gorm:"not null;unique"`
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key"`
	Name        string       `gorm:"not null;unique"`            // e.g. "admin", "recruiter"
	Permissions []Permission `gorm:"many2many:role_permissions"` // Many-to-Many relationship
}

type Event struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key"`
	Title         string         `gorm:"default:null"`
	EventMetaData datatypes.JSON `gorm:"type:json"`
	gorm.Model
}
