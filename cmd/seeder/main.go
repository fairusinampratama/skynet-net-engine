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

	// 2. Data to Seed
	name := "Main_Router"
	host := "tunnel.ebilling.id"
	port := 3724
	user := "skysky"
	pass := "skylineR34!@#"

	// 3. Upsert Logic (Insert if not exists)
	// We check by host to avoid duplicates
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM routers WHERE host = ?", host).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists > 0 {
		fmt.Printf("⚠️ Router %s already exists. Skipping.\n", host)
		return
	}

	_, err = db.Exec("INSERT INTO routers (name, host, port, username, password) VALUES (?, ?, ?, ?, ?)", 
		name, host, port, user, pass)
	if err != nil {
		log.Fatal("Failed to seed router:", err)
	}

	fmt.Printf("🎉 Successfully seeded router: %s@%s:%d\n", user, host, port)
}
