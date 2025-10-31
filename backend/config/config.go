package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// CORS
	AllowedOrigins []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Default JWT expiry: 7 days (7 * 24 hours = 168 hours)
	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "168"))

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")
	origins := strings.Split(allowedOrigins, ",")

	return &Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("ENV", "development"),
		DBDriver:       getEnv("DB_DRIVER", "postgres"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "8181"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "jemmy1303"),
		DBName:         getEnv("DB_NAME", "ganttpro_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours: jwtExpiryHours,
		AllowedOrigins: origins,
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
