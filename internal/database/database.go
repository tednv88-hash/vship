package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"vship/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error

	dsn := cfg.Database.DSN()

	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	isProd := appEnv == "prod" || appEnv == "production"
	logLevel := strings.ToLower(strings.TrimSpace(os.Getenv("GORM_LOG_LEVEL")))
	var mode logger.LogLevel
	switch logLevel {
	case "silent":
		mode = logger.Silent
	case "error":
		mode = logger.Error
	case "warn", "warning":
		mode = logger.Warn
	case "info":
		mode = logger.Info
	case "":
		if isProd {
			mode = logger.Warn
		} else {
			mode = logger.Info
		}
	default:
		if isProd {
			mode = logger.Warn
		} else {
			mode = logger.Info
		}
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(mode),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
