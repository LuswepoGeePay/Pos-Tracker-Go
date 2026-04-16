package main

import (
	"fmt"
	"log"
	"os"
	database "pos-master/config"
	"pos-master/routes"
	posservices "pos-master/services/pos_services"
	"pos-master/utils"
	"pos-master/utils/sentry"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize databaseuration from environment variables
	database.LoadEnv()

	err := utils.InitLogger("POS_MASTER.LOG")

	if err != nil {
		panic("Failed to initialize logger:" + fmt.Sprintf("error: %v", err))
	}

	database.InitDB()

	posservices.StartCronJobs()

	if err := sentry.InitSentry(); err != nil {
		panic("Sentry initialization failed:" + fmt.Sprintf("error: %v", err))
	}

	r := gin.Default()

	r.Use(sentrygin.New(sentrygin.Options{}))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true                                             // Allow all origins, or specify specific origins
	corsConfig.AllowMethods = []string{"GET", "POST", "DELETE", "PUT", "PATCH"}   // Allow specific HTTP methods
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // Allow specific headers

	r.Use(cors.New(corsConfig))

	routes.SetupRoutes(r)

	serverAddr := ":" + "8050"
	log.Printf("Starting server at %s (Environment: %s)", serverAddr, os.Getenv("ENVIRONMENT"))
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
