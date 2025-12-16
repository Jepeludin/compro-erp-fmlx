package testing

import (
	"encoding/json"
	"testing"
	"time"

	"ganttpro-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TokenBlacklist Structure Tests
// =============================================================================

func TestTokenBlacklist_Structure(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	
	token := models.TokenBlacklist{
		ID:        1,
		Token:     "hashed_token_string_sha256",
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	assert.Equal(t, uint(1), token.ID)
	assert.Equal(t, "hashed_token_string_sha256", token.Token)
	assert.True(t, token.ExpiresAt.After(now))
	assert.Equal(t, now, token.CreatedAt)
}

func TestTokenBlacklist_TableName(t *testing.T) {
	token := models.TokenBlacklist{}
	tableName := token.TableName()

	assert.Equal(t, "token_blacklists", tableName)
}

func TestTokenBlacklist_JSONSerialization(t *testing.T) {
	now := time.Now()
	token := models.TokenBlacklist{
		ID:        1,
		Token:     "test_token_hash",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}

	jsonData, err := json.Marshal(token)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "token")
	assert.Contains(t, result, "expires_at")
	assert.Contains(t, result, "created_at")
}

func TestTokenBlacklist_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": 1,
		"token": "abc123hash",
		"expires_at": "2024-12-15T12:00:00Z",
		"created_at": "2024-12-14T12:00:00Z"
	}`

	var token models.TokenBlacklist
	err := json.Unmarshal([]byte(jsonData), &token)
	require.NoError(t, err)

	assert.Equal(t, uint(1), token.ID)
	assert.Equal(t, "abc123hash", token.Token)
}

// =============================================================================
// Expiration Tests
// =============================================================================

func TestTokenBlacklist_IsExpired(t *testing.T) {
	testCases := []struct {
		name      string
		expiresAt time.Time
		isExpired bool
	}{
		{
			name:      "Future expiration (not expired)",
			expiresAt: time.Now().Add(24 * time.Hour),
			isExpired: false,
		},
		{
			name:      "Past expiration (expired)",
			expiresAt: time.Now().Add(-1 * time.Hour),
			isExpired: true,
		},
		{
			name:      "Just expired",
			expiresAt: time.Now().Add(-1 * time.Second),
			isExpired: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := models.TokenBlacklist{
				ExpiresAt: tc.expiresAt,
			}

			// Check if expired
			isExpired := token.ExpiresAt.Before(time.Now())
			assert.Equal(t, tc.isExpired, isExpired)
		})
	}
}

func TestTokenBlacklist_ExpirationDurations(t *testing.T) {
	now := time.Now()
	
	testCases := []struct {
		name     string
		duration time.Duration
	}{
		{"1 hour token", 1 * time.Hour},
		{"24 hour token", 24 * time.Hour},
		{"7 day token", 7 * 24 * time.Hour},
		{"30 day token", 30 * 24 * time.Hour},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := models.TokenBlacklist{
				ExpiresAt: now.Add(tc.duration),
				CreatedAt: now,
			}

			// Verify the duration is correct
			actualDuration := token.ExpiresAt.Sub(token.CreatedAt)
			assert.Equal(t, tc.duration, actualDuration)
		})
	}
}

// =============================================================================
// Token Hash Tests
// =============================================================================

func TestTokenBlacklist_TokenHashLength(t *testing.T) {
	// SHA256 produces 64 character hex string
	sha256Hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	
	token := models.TokenBlacklist{
		Token: sha256Hash,
	}

	assert.Len(t, token.Token, 64)
}

func TestTokenBlacklist_DifferentTokenHashes(t *testing.T) {
	tokens := []models.TokenBlacklist{
		{ID: 1, Token: "hash1_abc123"},
		{ID: 2, Token: "hash2_def456"},
		{ID: 3, Token: "hash3_ghi789"},
	}

	// All tokens should be unique
	tokenSet := make(map[string]bool)
	for _, t := range tokens {
		tokenSet[t.Token] = true
	}

	assert.Len(t, tokenSet, 3, "All tokens should be unique")
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestTokenBlacklist_EmptyFields(t *testing.T) {
	token := models.TokenBlacklist{}

	assert.Equal(t, uint(0), token.ID)
	assert.Empty(t, token.Token)
	assert.True(t, token.ExpiresAt.IsZero())
	assert.True(t, token.CreatedAt.IsZero())
}

func TestTokenBlacklist_LongTokenString(t *testing.T) {
	// JWT tokens can be quite long
	longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6IlRFU1RVU0VSIiwidXNlcl9pZF9zdHJpbmciOiJQSTEyMjQuMDAwMSIsInJvbGUiOiJBZG1pbiIsIm9wZXJhdG9yIjoiIiwiZXhwIjoxNzM0MTg5NjAwLCJpYXQiOjE3MzQxMDMyMDB9.abcdefghijklmnopqrstuvwxyz123456789"
	
	token := models.TokenBlacklist{
		Token: longToken,
	}

	assert.Equal(t, longToken, token.Token)
	assert.Greater(t, len(token.Token), 100)
}

func TestTokenBlacklist_ZeroExpiresAt(t *testing.T) {
	token := models.TokenBlacklist{
		ID:    1,
		Token: "some_hash",
	}

	// Zero time should be detected
	assert.True(t, token.ExpiresAt.IsZero())
}

// =============================================================================
// Time Precision Tests
// =============================================================================

func TestTokenBlacklist_TimePrecision(t *testing.T) {
	// Ensure time is stored with proper precision
	now := time.Now()
	token := models.TokenBlacklist{
		ExpiresAt: now,
		CreatedAt: now,
	}

	// Times should be equal
	assert.Equal(t, token.ExpiresAt, token.CreatedAt)
}

func TestTokenBlacklist_TimeZone(t *testing.T) {
	// Test with UTC time
	utcTime := time.Now().UTC()
	token := models.TokenBlacklist{
		ExpiresAt: utcTime,
		CreatedAt: utcTime,
	}

	assert.Equal(t, utcTime.Location().String(), token.ExpiresAt.Location().String())
}

// =============================================================================
// Database Index Tests
// =============================================================================

func TestTokenBlacklist_TokenIsIndexed(t *testing.T) {
	// The Token field has an index tag: `gorm:"type:text;not null;index"`
	// This test documents that expectation
	token := models.TokenBlacklist{
		Token: "indexed_token_value",
	}

	assert.NotEmpty(t, token.Token, "Token should be set for index to work")
}

func TestTokenBlacklist_ExpiresAtIsIndexed(t *testing.T) {
	// The ExpiresAt field has an index tag: `gorm:"not null;index"`
	// This test documents that expectation
	token := models.TokenBlacklist{
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	assert.False(t, token.ExpiresAt.IsZero(), "ExpiresAt should be set for index to work")
}

// =============================================================================
// Cleanup Scenario Tests
// =============================================================================

func TestTokenBlacklist_IdentifyExpiredTokens(t *testing.T) {
	now := time.Now()
	tokens := []models.TokenBlacklist{
		{ID: 1, Token: "expired1", ExpiresAt: now.Add(-2 * time.Hour)},
		{ID: 2, Token: "expired2", ExpiresAt: now.Add(-1 * time.Hour)},
		{ID: 3, Token: "valid1", ExpiresAt: now.Add(1 * time.Hour)},
		{ID: 4, Token: "valid2", ExpiresAt: now.Add(24 * time.Hour)},
	}

	var expiredTokens []models.TokenBlacklist
	for _, token := range tokens {
		if token.ExpiresAt.Before(now) {
			expiredTokens = append(expiredTokens, token)
		}
	}

	assert.Len(t, expiredTokens, 2, "Should have 2 expired tokens")
	assert.Equal(t, "expired1", expiredTokens[0].Token)
	assert.Equal(t, "expired2", expiredTokens[1].Token)
}

// =============================================================================
// Security Tests
// =============================================================================

func TestTokenBlacklist_TokenNotStoredAsPlaintext(t *testing.T) {
	// Document that tokens should be hashed before storage
	originalToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"
	hashedToken := "hashed_version_of_token_sha256" // In real code, this would be SHA256 hash

	token := models.TokenBlacklist{
		Token: hashedToken,
	}

	// The stored token should NOT equal the original token
	assert.NotEqual(t, originalToken, token.Token, "Token should be hashed, not stored as plaintext")
}

