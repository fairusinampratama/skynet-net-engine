package main

import (
	"skynet-net-engine-api/internal/database"
	"skynet-net-engine-api/internal/core"
	"skynet-net-engine-api/internal/api"
	"skynet-net-engine-api/pkg/logger"
	"go.uber.org/zap"
	_ "skynet-net-engine-api/docs" // Required for Swagger
)

// @title           NetEngine API
// @version         1.0
// @description     Middleware for managing Mikrotik Routers
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @securityDefinitions.apikey AppKey
// @in header
// @name X-App-Key
func main() {
	// 1. Initialize Logger
	logger.Init()
	logger.Info("Starting NetEngine...", zap.String("version", "1.0.0"))

	// 2. Initialize Database
	database.Init()
	defer database.DB.Close()

	// 3. Initialize Worker Pool (Spawns 1 Goroutine per Router)
	core.InitPool()

	// 4. Start HTTP API
	go api.Start(":8080")

	// Block forever
	select {}
}
