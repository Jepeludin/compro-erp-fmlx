package testing

import (
	"testing"

	"ganttpro-backend/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// HashPassword Tests
// =============================================================================

func TestHashPassword_Success(t *testing.T) {
	// Arrange
	password := "securePassword123"

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err, "HashPassword should not return an error")
	assert.NotEmpty(t, hash, "Hash should not be empty")
	assert.NotEqual(t, password, hash, "Hash should be different from original password")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	// Arrange
	password := ""

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err, "HashPassword should handle empty password")
	assert.NotEmpty(t, hash, "Hash should not be empty even for empty password")
}

func TestHashPassword_LongPassword(t *testing.T) {
	// Arrange - bcrypt has a max length of 72 bytes
	password := "ThisIsAVeryLongPasswordThatExceedsTheNormalLengthButShouldStillWork123456789"

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err, "HashPassword should handle long passwords")
	assert.NotEmpty(t, hash, "Hash should not be empty")
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		password string
	}{
		{"WithSymbols", "P@ssw0rd!#$%"},
		{"WithUnicode", "密码123"},
		{"WithSpaces", "pass word with spaces"},
		{"WithNewline", "password\nwith\nnewlines"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			hash, err := utils.HashPassword(tc.password)

			// Assert
			require.NoError(t, err, "HashPassword should handle special characters")
			assert.NotEmpty(t, hash, "Hash should not be empty")
		})
	}
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	// Arrange
	password := "samePassword123"

	// Act
	hash1, err1 := utils.HashPassword(password)
	hash2, err2 := utils.HashPassword(password)

	// Assert
	require.NoError(t, err1, "First hash should not return an error")
	require.NoError(t, err2, "Second hash should not return an error")
	assert.NotEqual(t, hash1, hash2, "Same password should produce different hashes (due to salt)")
}

// =============================================================================
// CheckPasswordHash Tests
// =============================================================================

func TestCheckPasswordHash_ValidPassword(t *testing.T) {
	// Arrange
	password := "correctPassword123"
	hash, _ := utils.HashPassword(password)

	// Act
	result := utils.CheckPasswordHash(password, hash)

	// Assert
	assert.True(t, result, "CheckPasswordHash should return true for correct password")
}

func TestCheckPasswordHash_InvalidPassword(t *testing.T) {
	// Arrange
	password := "correctPassword123"
	wrongPassword := "wrongPassword456"
	hash, _ := utils.HashPassword(password)

	// Act
	result := utils.CheckPasswordHash(wrongPassword, hash)

	// Assert
	assert.False(t, result, "CheckPasswordHash should return false for incorrect password")
}

func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	// Arrange
	password := ""
	hash, _ := utils.HashPassword(password)

	// Act
	result := utils.CheckPasswordHash(password, hash)

	// Assert
	assert.True(t, result, "CheckPasswordHash should work with empty password")
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	// Arrange
	password := "somePassword"
	invalidHash := "notAValidBcryptHash"

	// Act
	result := utils.CheckPasswordHash(password, invalidHash)

	// Assert
	assert.False(t, result, "CheckPasswordHash should return false for invalid hash")
}

func TestCheckPasswordHash_EmptyHash(t *testing.T) {
	// Arrange
	password := "somePassword"
	emptyHash := ""

	// Act
	result := utils.CheckPasswordHash(password, emptyHash)

	// Assert
	assert.False(t, result, "CheckPasswordHash should return false for empty hash")
}

func TestCheckPasswordHash_CaseSensitive(t *testing.T) {
	// Arrange
	password := "CaseSensitivePassword"
	hash, _ := utils.HashPassword(password)

	// Act & Assert
	assert.True(t, utils.CheckPasswordHash("CaseSensitivePassword", hash), "Exact match should succeed")
	assert.False(t, utils.CheckPasswordHash("casesensitivepassword", hash), "Lowercase should fail")
	assert.False(t, utils.CheckPasswordHash("CASESENSITIVEPASSWORD", hash), "Uppercase should fail")
}

// =============================================================================
// Integration Tests (Hash and Check together)
// =============================================================================

func TestPasswordHashingRoundTrip(t *testing.T) {
	// Table-driven test for various password scenarios
	testCases := []struct {
		name     string
		password string
	}{
		{"SimplePassword", "password123"},
		{"ComplexPassword", "C0mpl3x!P@ssw0rd#2024"},
		{"NumericPassword", "1234567890"},
		{"ShortPassword", "abc"},
		{"PasswordWithSpaces", "my secret password"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange & Act
			hash, err := utils.HashPassword(tc.password)

			// Assert
			require.NoError(t, err, "Hashing should not fail")
			assert.True(t, utils.CheckPasswordHash(tc.password, hash), "Password should match its hash")
			assert.False(t, utils.CheckPasswordHash(tc.password+"wrong", hash), "Wrong password should not match")
		})
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"
	for i := 0; i < b.N; i++ {
		utils.HashPassword(password)
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := utils.HashPassword(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.CheckPasswordHash(password, hash)
	}
}
