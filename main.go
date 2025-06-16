package main

import (
	"fmt"
	"log"
	"pos-master/config"
	"pos-master/routes"
	"pos-master/utils"
	"pos-master/utils/sentry"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	err := utils.InitLogger("pos_master.log")

	if err != nil {
		panic("Failed to initialize logger:" + fmt.Sprintf("error: %v", err))
	}

	config.InitDB()
	if err := sentry.InitSentry(); err != nil {
		panic("Sentry initialization failed:" + fmt.Sprintf("error: %v", err))
	}

	r := gin.Default()

	r.Use(sentrygin.New(sentrygin.Options{}))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                             // Allow all origins, or specify specific origins
	config.AllowMethods = []string{"GET", "POST", "DELETE", "PUT", "PATCH"}   // Allow specific HTTP methods
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // Allow specific headers

	r.Use(cors.New(config))

	routes.SetupRoutes(r)

	// certFile := os.Getenv("CERT_FILE")
	// keyFile := os.Getenv("KEY_FILE")

	// if certFile != "" && keyFile != "" {
	// 	// Verify certificate files exist
	// 	if _, err := os.Stat(certFile); err != nil {
	// 		log.Printf("Warning: Certificate file not found: %v", err)
	// 		log.Println("Falling back to HTTP")
	// 	} else if _, err := os.Stat(keyFile); err != nil {
	// 		log.Printf("Warning: Key file not found: %v", err)
	// 		log.Println("Falling back to HTTP")
	// 	} else {
	// 		// Start HTTPS server
	// 		log.Println("Starting HTTPS server on port 7443")
	// 		if err := r.RunTLS(":7443", certFile, keyFile); err != nil {
	// 			log.Fatalf("Failed to start HTTPS server: %v", err)
	// 		}
	// 		return
	// 	}
	// }

	log.Println("Starting server at 8050")
	if err := r.Run(":8050"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
