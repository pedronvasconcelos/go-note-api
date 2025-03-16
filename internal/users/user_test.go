package users

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// Arrange
	firstName := "John"
	lastName := "Doe"
	email := "john.doe@example.com"

	// Act
	user := newUser(firstName, lastName, email)

	// Assert
	assert.NotEqual(t, uuid.Nil, user.ID, "User ID should not be nil")
	assert.Equal(t, firstName, user.FirstName, "FirstName should match")
	assert.Equal(t, lastName, user.LastName, "LastName should match")
	assert.Equal(t, email, user.Email, "Email should match")
	assert.NotZero(t, user.CreatedAt, "CreatedAt should be set")
	assert.NotZero(t, user.UpdatedAt, "UpdatedAt should be set")
	assert.Nil(t, user.PremiumUntil, "PremiumUntil should be nil")
	assert.True(t, user.isActive, "isActive should be true")
}

func TestIsPremium(t *testing.T) {
	tests := []struct {
		name           string
		premiumUntil   *time.Time
		expectedResult bool
	}{
		{
			name:           "User with no premium",
			premiumUntil:   nil,
			expectedResult: false,
		},
		{
			name:           "User with expired premium",
			premiumUntil:   timePtr(time.Now().Add(-24 * time.Hour)),
			expectedResult: false,
		},
		{
			name:           "User with active premium",
			premiumUntil:   timePtr(time.Now().Add(24 * time.Hour)),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			user := &User{PremiumUntil: tt.premiumUntil}

			// Act
			result := user.isPremium()

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestActivatePremium30Days(t *testing.T) {
	// Arrange
	user := newUser("Jane", "Doe", "jane.doe@example.com")
	initialUpdateTime := user.UpdatedAt

	// Make sure we wait a tiny bit to ensure UpdatedAt changes
	time.Sleep(1 * time.Millisecond)

	// Act
	user.activatePremium30Days()

	// Assert
	assert.NotNil(t, user.PremiumUntil, "PremiumUntil should not be nil")

	expectedTime := time.Now().Add(30 * 24 * time.Hour)
	timeDiff := expectedTime.Sub(*user.PremiumUntil)
	assert.LessOrEqual(t, timeDiff.Abs(), 2*time.Second, "PremiumUntil should be approximately 30 days from now")

	assert.True(t, user.UpdatedAt.After(initialUpdateTime), "UpdatedAt should be updated")
	assert.True(t, user.isPremium(), "User should be premium after activation")
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
