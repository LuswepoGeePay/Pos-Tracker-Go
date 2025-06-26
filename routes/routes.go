package routes

import (
	"pos-master/controllers/apps"
	"pos-master/controllers/auth"
	"pos-master/controllers/dashboard"
	"pos-master/controllers/events"
	locationhistory "pos-master/controllers/location_history"
	posdevices "pos-master/controllers/pos_devices"
	"pos-master/controllers/users"
	"pos-master/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/v1/create-user", auth.RegisterHandler)
	r.POST("/v1/login", auth.LoginHandler)

	auth := r.Group("/v1")
	auth.Use(middleware.AuthMiddleware())

	//users

	auth.POST("/users/get", users.GetUsersHandler)
	auth.POST("/user/update", users.EditUserHandler)

	//reports
	// auth.POST("/report/devices/generate")

	//dashboard
	auth.GET("/dashboard/tiles/get", dashboard.GetTileInfoHandler)
	auth.GET("/dashboard/pie/get", dashboard.GetPieChartDataHandler)
	auth.GET("/dashboard/bar/get", dashboard.GetLineChartHandler)
	auth.POST("/dashboard/events/get", events.GetEventsHandler)
	//Pos devices
	r.POST("/v1/pos/register", posdevices.RegisterPosDeviceHandler)
	auth.POST("/pos/devices/get", posdevices.GetPosDevicesHandler)
	auth.POST("/pos/device/update", posdevices.EditDeviceHandler)
	auth.DELETE("/pos/device/:id", posdevices.DeleteDeviceHandler)
	r.POST("/v1/pos/device/heartbeat", posdevices.HeartBeatHandler)

	//Apps
	auth.POST("/app/register", apps.RegisterAppHandler)
	auth.POST("/apps/get", apps.GetAppsHandler)
	r.POST("/v1/app/update", apps.CheckAppUpdate)
	auth.POST("/app/info/update", apps.EditAppHandler)
	auth.DELETE("/app/:id", apps.DeleteAppHandler)

	// auth.DELETE("/app/delete/:id")

	//App versions
	auth.POST("/app/version/register", apps.RegisterNewAppVersionHandler)
	auth.POST("/app/versions/get", apps.GetAppVersionsHandler)
	auth.POST("/app/version/update", apps.EditAppVersionHandler)
	auth.DELETE("/app/version/:id", apps.DeleteAppVersionHandler)

	//location history
	r.POST("/v1/location/register", locationhistory.RegisterNewLocationHandler)
	auth.POST("/locations/get", locationhistory.GetLocationsHandler)
	// auth.DELETE("/pos/device/id", posdevices.DeleteDeviceHandler)

}
