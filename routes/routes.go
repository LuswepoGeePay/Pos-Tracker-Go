package routes

import (
	"pos-master/controllers/apps"
	"pos-master/controllers/auth"
	locationhistory "pos-master/controllers/location_history"
	posdevices "pos-master/controllers/pos_devices"
	"pos-master/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/create-user", auth.RegisterHandler)
	r.POST("/login", auth.LoginHandler)

	auth := r.Group("/v1")
	auth.Use(middleware.AuthMiddleware())

	//Pos devices
	r.POST("/v1/pos/register", posdevices.RegisterPosDeviceHandler)
	auth.POST("/v1/pos/devices/get", posdevices.GetPosDevicesHandler)

	//Apps
	auth.POST("/app/register", apps.RegisterAppHandler)
	auth.POST("/apps/get", apps.GetAppsHandler)

	//App versions
	auth.POST("/app/version/register", apps.RegisterNewAppVersionHandler)
	auth.POST("/app/versions/get", apps.GetAppVersionsHandler)
	//location history
	auth.POST("/location/register", locationhistory.RegisterNewLocationHandler)
	auth.POST("/locations/get", locationhistory.GetLocationsHandler)

}
