package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/wafi04/chatting-app/config/database"
	"github.com/wafi04/chatting-app/config/env"
	"github.com/wafi04/chatting-app/services/gateway"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
)

func main() {
	logs := logger.NewLogger()
	log := logger.NewLogger()

	db, err := database.NewDB(env.LoadEnv("DB_URL"))
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize database: %v", err)
		return
	}

	defer db.Close()

	health := db.Health()
	log.Log(logger.InfoLevel, "Database health: %v", health["status"])

	logs.Info("Staring Server gateway ")

	router := gateway.SetUpRoutes(db)

	if err := router.Run(":8080"); err != nil {
		logs.Log(logger.ErrorLevel, "Failed to start server: %s", err)
	}
}
