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

	//Apps
	auth.POST("/app/register", apps.RegisterAppHandler)

	//App versions
	// auth.POST("/v1/app/version/register", apps.)

	//location history
	auth.POST("/location/register", locationhistory.RegisterNewLocationHandler)
}
