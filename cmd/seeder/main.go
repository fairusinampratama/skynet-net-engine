package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 1. Connection logic (Simplified version of internal/database for the script)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "fairusinampratama@tcp(127.0.0.1:3306)/netengine?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	fmt.Println("✅ Connected to Database")

	// 2. Data to Seed (From User Image)
	// Two credential sets based on user feedback:
	// - Most routers: skynet / sky123!@#
	// - Routers 1,6,8,11,13,15: skysky / skylineR34!@#
	
	routers := []struct {
		Name string
		Host string
		Port int
		User string
		Pass string
	}{
		{"Randuagung-CCR", "tunnel.ebilling.id", 3724, "skysky", "skylineR34!@#"},
		{"Skynet Srigading", "tunnel.ebilling.id", 1973, "skysky", "skylineR34!@#"},
		{"Skynet Arjosari", "tunnel.ebilling.id", 3718, "skynet", "sky123!@#"},
		{"Skynet Krian", "103.156.128.34", 8777, "skynet", "sky123!@#"},
		{"Skynet-Rest-Area-Karang-Ploso", "tunnel.ebilling.id", 3625, "skynet", "sky123!@#"},
		{"Skynet Lawang", "tunnel.ebilling.id", 3499, "skysky", "skylineR34!@#"},
		{"Skynet Kunci", "tunnel.ebilling.id", 3496, "skynet", "sky123!@#"},
		{"Skynet Purwosari - Purwodadi", "tunnel.ebilling.id", 16980, "skysky", "skylineR34!@#"},
		{"Skynet Tutur", "tunnel.ebilling.id", 8200, "skysky", "skylineR34!@#"},
		{"Skynet Bukit Sentul", "tunnel.ebilling.id", 9529, "skynet", "sky123!@#"},
		{"Skynet Bantaran", "tunnel.ebilling.id", 14939, "skysky", "skylineR34!@#"},
		{"Skynet Kasin", "tunnel.ebilling.id", 2734, "skynet", "sky123!@#"},
		{"Skynet Tasikmadu", "tunnel.ebilling.id", 15515, "skysky", "skylineR34!@#"},
		{"Skynet Kendit", "tunnel.ebilling.id", 16295, "skynet", "sky123!@#"},
		{"Skynet Bumiayu", "tunnel2.ebilling.id", 23506, "skysky", "skylineR34!@#"},
	}

	fmt.Printf("🌱 Seeding %d routers...\n", len(routers))

	for _, r := range routers {
		// Check duplicates
		var exists int
		err = db.QueryRow("SELECT COUNT(*) FROM routers WHERE host = ? AND port = ?", r.Host, r.Port).Scan(&exists)
		if err != nil {
			log.Printf("❌ Checking %s: %v\n", r.Name, err)
			continue
		}

		if exists > 0 {
			fmt.Printf("⚠️  Skipping %s (Already exists)\n", r.Name)
			continue
		}

		_, err = db.Exec("INSERT INTO routers (name, host, port, username, password) VALUES (?, ?, ?, ?, ?)", 
			r.Name, r.Host, r.Port, r.User, r.Pass)
		
		if err != nil {
			log.Printf("❌ Failed to seed %s: %v\n", r.Name, err)
		} else {
			fmt.Printf("✅ Seeded: %s\n", r.Name)
		}
	}
	
	fmt.Println("🎉 Seeding Complete!")
}
