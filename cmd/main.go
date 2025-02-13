package main

import (
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wafi04/chatting-app/config/database"
	"github.com/wafi04/chatting-app/config/env"
	"github.com/wafi04/chatting-app/services/gateway"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
)

func validateCloudinaryConfig() (cloudName, apiKey, apiSecret string, err error) {
	cloudName = strings.Trim(env.LoadEnv("CLOUDINARY_CLOUD_NAME"), "\"")
	apiKey = strings.Trim(env.LoadEnv("CLOUDINARY_API_KEY"), "\"")
	apiSecret = strings.Trim(env.LoadEnv("CLOUDINARY_API_SECRET"), "\"")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", "", "", fmt.Errorf("missing required Cloudinary configuration")
	}

	return cloudName, apiKey, apiSecret, nil
}

func main() {
	logs := logger.NewLogger()
	log := logger.NewLogger()

	// Database initialization
	dbURL := env.LoadEnv("DB_URL")
	if dbURL == "" {
		log.Log(logger.ErrorLevel, "DB_URL environment variable is not set")
		return
	}

	db, err := database.NewDB(dbURL)
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize database: %v", err)
		return
	}
	defer db.Close()

	mongo, err := database.ConnectMongoDB(log)
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize mongo database: %v", err)
		return
	}
	defer mongo.Close()

	// Cloudinary configuration validation
	cloudName, apiKey, apiSecret, err := validateCloudinaryConfig()
	if err != nil {
		log.Log(logger.ErrorLevel, "Cloudinary configuration error: %v", err)
		return
	}

	// Debug logging with length and trimmed values
	log.Log(logger.InfoLevel, "Cloudinary Cloud Name %s :", env.LoadEnv("MONGO_URL"))

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize Cloudinary: %v", err)
		return
	}

	health := db.Health()
	log.Log(logger.InfoLevel, "Database health: %v", health["status"])

	logs.Info("Starting Server gateway")

	router := gateway.SetUpRoutes(db, mongo.Client, cld)
	port := env.LoadEnv("PORT")
	if port == "" {
		port = ":8080" // default port if not set
	}

	if err := router.Run(port); err != nil {
		logs.Log(logger.ErrorLevel, "Failed to start server: %s", err)
	}
}
