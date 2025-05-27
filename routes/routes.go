package routes

import (
	"pos-master/controllers/auth"
	"pos-master/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/create-user", auth.RegisterHandler)
	r.POST("/login", auth.LoginHandler)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
}
