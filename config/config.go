package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database    DatabaseConfig
	Server      ServerConfig
	JWT         JWTConfig
	GoogleOAuth GoogleOAuthConfig
	SMTP        SMTPConfig
	CompanyName string
	AppName     string
	BaseURL     string
}

type GoogleOAuthConfig struct {
	Enabled  bool
	ClientID string
}

type SMTPConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	FromEmail string
	FromName  string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Host string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			URL:      getEnv("DATABASE_URL", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "vship"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: firstEnv("SERVER_PORT", "PORT", "3002"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "vship-secret-key-change-in-production"),
		},
		GoogleOAuth: GoogleOAuthConfig{
			Enabled:  getEnv("GOOGLE_OAUTH_ENABLED", "") == "true",
			ClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", ""),
			Port:      getEnv("SMTP_PORT", "587"),
			User:      getEnv("SMTP_USER", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			FromEmail: getEnv("SMTP_FROM_EMAIL", ""),
			FromName:  getEnv("SMTP_FROM_NAME", "vShip"),
		},
		CompanyName: getEnv("COMPANY_NAME", "V-sys Limited"),
		AppName:     getEnv("APP_NAME", "vShip"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:3002"),
	}
}

func (c *DatabaseConfig) DSN() string {
	if c.URL != "" {
		return c.URL
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func firstEnv(keys ...string) string {
	for _, key := range keys[:len(keys)-1] {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return keys[len(keys)-1]
}
