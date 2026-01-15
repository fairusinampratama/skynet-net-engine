package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"skynet-net-engine-api/pkg/logger"
	"go.uber.org/zap"
)

var DB *sql.DB

func Init() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Default fallback for development
		// User provided: fairusinampratama, NO pass
		dsn = "fairusinampratama@tcp(127.0.0.1:3306)/netengine?parseTime=true"
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to open database connection")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	if err := DB.Ping(); err != nil {
		logger.Warn("Failed to ping database - Running in Offline/Mock Mode", zap.Error(err))
		// We don't Fatal here to allow testing the API without a DB
		// In production, this should probably fatal.
	} else {
		logger.Info("Database connected successfully")
	}
}
