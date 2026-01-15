package main

import (
	"log"
	"skynet-net-engine-api/internal/database"
	"skynet-net-engine-api/internal/mikrotik"
	"skynet-net-engine-api/pkg/logger"
	
	"go.uber.org/zap"
)

// SyncUsersFromMikrotik fetches all PPPoE secrets from MikroTik and syncs to database
func main() {
	logger.Init()
	database.Init()

	log.Println("🔄 Starting user sync from MikroTik...")

	// Get all routers
	routers, err := database.GetAllRouters()
	if err != nil {
		log.Fatalf("Failed to fetch routers: %v", err)
	}

	totalSynced := 0
	for _, router := range routers {
		log.Printf("📡 Syncing users from router: %s (%s)", router.Name, router.Host)

		// Connect to MikroTik
		client, err := mikrotik.NewClient(router)
		if err != nil {
			log.Printf("❌ Failed to connect to %s: %v", router.Host, err)
			continue
		}

		// Fetch all PPPoE secrets
		secrets, err := client.GetAllSecrets()
		if err != nil {
			log.Printf("❌ Failed to fetch secrets from %s: %v", router.Host, err)
			client.Close()
			continue
		}

		// Insert into database
		for _, secret := range secrets {
			err := database.UpsertUser(secret.Name, router.ID, secret.Profile)
			if err != nil {
				logger.Error("Failed to insert user", zap.String("user", secret.Name), zap.Error(err))
			} else {
				totalSynced++
			}
		}

		client.Close()
		log.Printf("✅ Synced %d users from %s", len(secrets), router.Name)
	}

	log.Printf("🎉 Total users synced: %d", totalSynced)
}
