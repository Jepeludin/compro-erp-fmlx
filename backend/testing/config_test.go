package testing

import (
	"os"
	"testing"
	"ganttpro-backend/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Config Structure Tests
// =============================================================================

func TestConfig_Structure(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		Environment:    "development",
		DBDriver:       "postgres",
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "testuser",
		DBPassword:     "testpass",
		DBName:         "testdb",
		DBSSLMode:      "disable",
		JWTSecret:      "secret",
		JWTExpiryHours: 24,
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "postgres", cfg.DBDriver)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "testpass", cfg.DBPassword)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, 24, cfg.JWTExpiryHours)
	assert.Len(t, cfg.AllowedOrigins, 1)
}

// =============================================================================
// LoadConfig Tests
// =============================================================================

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear all environment variables
	envVars := []string{
		"PORT", "ENV", "DB_DRIVER", "DB_HOST", "DB_PORT",
		"DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"JWT_SECRET", "JWT_EXPIRY_HOURS", "ALLOWED_ORIGINS",
	}
	
	// Store original values
	originalValues := make(map[string]string)
	for _, key := range envVars {
		originalValues[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	
	// Restore after test
	defer func() {
		for key, value := range originalValues {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	cfg := config.LoadConfig()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "postgres", cfg.DBDriver)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, 168, cfg.JWTExpiryHours) // 7 days default
	assert.Contains(t, cfg.AllowedOrigins, "http://localhost:5173")
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Store original values
	originalValues := map[string]string{
		"PORT":             os.Getenv("PORT"),
		"ENV":              os.Getenv("ENV"),
		"DB_HOST":          os.Getenv("DB_HOST"),
		"DB_PORT":          os.Getenv("DB_PORT"),
		"JWT_SECRET":       os.Getenv("JWT_SECRET"),
		"JWT_EXPIRY_HOURS": os.Getenv("JWT_EXPIRY_HOURS"),
		"ALLOWED_ORIGINS":  os.Getenv("ALLOWED_ORIGINS"),
	}

	// Restore after test
	defer func() {
		for key, value := range originalValues {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	// Set test environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("ENV", "production")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("JWT_SECRET", "my-super-secret")
	os.Setenv("JWT_EXPIRY_HOURS", "48")
	os.Setenv("ALLOWED_ORIGINS", "http://example.com,http://app.example.com")

	cfg := config.LoadConfig()

	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "my-super-secret", cfg.JWTSecret)
	assert.Equal(t, 48, cfg.JWTExpiryHours)
	assert.Len(t, cfg.AllowedOrigins, 2)
	assert.Contains(t, cfg.AllowedOrigins, "http://example.com")
	assert.Contains(t, cfg.AllowedOrigins, "http://app.example.com")
}

func TestLoadConfig_InvalidJWTExpiryHours(t *testing.T) {
	// Store and clear JWT_EXPIRY_HOURS
	originalValue := os.Getenv("JWT_EXPIRY_HOURS")
	defer func() {
		if originalValue != "" {
			os.Setenv("JWT_EXPIRY_HOURS", originalValue)
		} else {
			os.Unsetenv("JWT_EXPIRY_HOURS")
		}
	}()

	// Set invalid value
	os.Setenv("JWT_EXPIRY_HOURS", "invalid")

	cfg := config.LoadConfig()

	// Should default to 0 when parsing fails
	assert.Equal(t, 0, cfg.JWTExpiryHours)
}

func TestLoadConfig_MultipleOrigins(t *testing.T) {
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	defer func() {
		if originalValue != "" {
			os.Setenv("ALLOWED_ORIGINS", originalValue)
		} else {
			os.Unsetenv("ALLOWED_ORIGINS")
		}
	}()

	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173,http://example.com")

	cfg := config.LoadConfig()

	assert.Len(t, cfg.AllowedOrigins, 3)
}

// =============================================================================
// getEnv Tests
// =============================================================================

func TestGetEnv_ReturnsValue(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := os.Getenv("TEST_VAR")

	assert.Equal(t, "test_value", result)
}

func TestGetEnv_ReturnsDefault(t *testing.T) {
	// Ensure variable doesn't exist
	os.Unsetenv("NONEXISTENT_VAR")

	result := os.Getenv("NONEXISTENT_VAR")

	assert.Equal(t, "default_value", result)
}

func TestGetEnv_EmptyValueReturnsDefault(t *testing.T) {
	// Set empty value
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result := os.Getenv("EMPTY_VAR")

	assert.Equal(t, "default_value", result)
}

// =============================================================================
// Config Validation Tests
// =============================================================================

func TestConfig_JWTSecretNotEmpty(t *testing.T) {
	cfg := config.LoadConfig()

	require.NotEmpty(t, cfg.JWTSecret, "JWT secret should have a default value")
}

func TestConfig_PortIsValid(t *testing.T) {
	cfg := config.LoadConfig()

	require.NotEmpty(t, cfg.Port, "Port should have a default value")
}

func TestConfig_DatabaseSettings(t *testing.T) {
	cfg := config.LoadConfig()

	assert.NotEmpty(t, cfg.DBDriver)
	assert.NotEmpty(t, cfg.DBHost)
	assert.NotEmpty(t, cfg.DBPort)
	assert.NotEmpty(t, cfg.DBUser)
	assert.NotEmpty(t, cfg.DBName)
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestLoadConfig_EmptyAllowedOrigins(t *testing.T) {
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	defer func() {
		if originalValue != "" {
			os.Setenv("ALLOWED_ORIGINS", originalValue)
		} else {
			os.Unsetenv("ALLOWED_ORIGINS")
		}
	}()

	os.Setenv("ALLOWED_ORIGINS", "")

	cfg := config.LoadConfig()

	// Empty string gets split into single empty element, then falls back to default
	assert.NotEmpty(t, cfg.AllowedOrigins)
}

func TestLoadConfig_SingleOrigin(t *testing.T) {
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	defer func() {
		if originalValue != "" {
			os.Setenv("ALLOWED_ORIGINS", originalValue)
		} else {
			os.Unsetenv("ALLOWED_ORIGINS")
		}
	}()

	os.Setenv("ALLOWED_ORIGINS", "http://single-origin.com")

	cfg := config.LoadConfig()

	assert.Len(t, cfg.AllowedOrigins, 1)
	assert.Equal(t, "http://single-origin.com", cfg.AllowedOrigins[0])
}

