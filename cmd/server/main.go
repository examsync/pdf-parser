package main

import (
	"github.com/examsync/pdf-parser/utils/config"
	"github.com/examsync/pdf-parser/utils/database"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	// 1. Initialize Logrus Logger
	logger.InitLogger()
	logger.Log.Info("Starting PDF Parser Service...")

	// 2. Load Configuration via Viper
	logger.Log.Info("Loading configuration using Viper...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load application configuration: %v", err)
	}

	// Print Viper Loaded configuration
	PrintConfiguration(cfg)

	// 3. Initialize Database Client
	db := initDatabase(&cfg.Database)
	defer func() {
		if err := database.CloseDB(db); err != nil {
			logger.Log.Errorf("Error closing database connection: %v", err)
		} else {
			logger.Log.Info("Database connection closed gracefully")
		}
	}()

	// 4. Start HTTP Web Server
	startServer(cfg.Server.Port, db)
}

// initDatabase connects to the database via GORM and fails fast if the connection is invalid
	func initDatabase(cfg *config.DatabaseConfig) *gorm.DB {
	logger.Log.Info("Connecting to PostgreSQL database via GORM...")
	db, err := database.ConnectDB(cfg)
	if err != nil {
		logger.Log.Fatalf("Database connection setup failed: %v", err)
	}
	return db
}

// PrintConfiguration outputs loaded Viper settings using Logrus
func PrintConfiguration(cfg *config.Config) {
	logger.Log.WithFields(logrus.Fields{
		"server.port":      cfg.Server.Port,
		"database.host":    cfg.Database.Host,
		"database.port":    cfg.Database.Port,
		"database.user":    cfg.Database.User,
		"database.dbname":  cfg.Database.DBName,
		"database.sslmode": cfg.Database.SSLMode,
	}).Info("Viper configuration loaded successfully")
}
