package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDB is a simple mock implementation for testing
type MockDB struct {
	users        map[uuid.UUID]*User
	usersByEmail map[string]*User
	createError  error
	updateError  error
	readError    error
	existsError  error
}

func NewMockDB() *MockDB {
	return &MockDB{
		users:        make(map[uuid.UUID]*User),
		usersByEmail: make(map[string]*User),
	}
}

// WithContext implements the gorm-like interface
func (m *MockDB) WithContext(_ context.Context) *MockDB {
	return m
}

func (m *MockDB) Create(user *User) *MockDB {
	if m.createError != nil {
		return m
	}
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	return m
}

func (m *MockDB) Save(user *User) *MockDB {
	if m.updateError != nil {
		return m
	}
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	return m
}

func (m *MockDB) Where(query string, args ...interface{}) *MockDB {
	return m
}

func (m *MockDB) Model(value interface{}) *MockDB {
	return m
}

func (m *MockDB) Count(count *int64) *MockDB {
	if m.existsError != nil {
		return m
	}

	email := ""
	*count = 0
	for _, u := range m.users {
		if u.Email == email {
			*count = 1
			break
		}
	}
	return m
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *MockDB {
	if m.readError != nil {
		return m
	}

	_, ok := dest.(*User)
	if !ok {
		return m
	}
	return m
}

// Error mocks error handling
func (m *MockDB) Error() error {
	// Return the appropriate error based on the operation
	if m.createError != nil {
		return m.createError
	}
	if m.updateError != nil {
		return m.updateError
	}
	if m.readError != nil {
		return m.readError
	}
	if m.existsError != nil {
		return m.existsError
	}
	return nil
}

// mockRepository creates a mock implementation of Repository
type mockRepository struct {
	mockDB *MockDB
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		mockDB: NewMockDB(),
	}
}

func (m *mockRepository) Add(_ context.Context, user *User) error {
	return m.mockDB.Create(user).Error()
}

func (m *mockRepository) Update(_ context.Context, user *User) error {
	return m.mockDB.Save(user).Error()
}

func (m *mockRepository) GetByEmail(_ context.Context, email string) (*User, error) {
	if m.mockDB.readError != nil {
		return nil, m.mockDB.readError
	}
	return m.mockDB.usersByEmail[email], nil
}

func (m *mockRepository) GetByID(_ context.Context, id uuid.UUID) (*User, error) {
	if m.mockDB.readError != nil {
		return nil, m.mockDB.readError
	}
	return m.mockDB.users[id], nil
}

func (m *mockRepository) ExistsByEmail(_ context.Context, email string) (bool, error) {
	if m.mockDB.existsError != nil {
		return false, m.mockDB.existsError
	}
	_, exists := m.mockDB.usersByEmail[email]
	return exists, nil
}

// createTestUser creates a test user for testing
func createTestUser() *User {
	return &User{
		ID:           uuid.New(),
		FirstName:    "Test",
		LastName:     "User",
		Email:        "test@example.com",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PremiumUntil: nil,
		isActive:     true,
	}
}

func TestRepositoryAdd(t *testing.T) {
	// Setup
	mockRepo := NewMockRepository()
	ctx := context.Background()

	// Create test user
	user := createTestUser()

	// Test Add
	err := mockRepo.Add(ctx, user)
	assert.NoError(t, err)

	// Verify user was added
	savedUser, err := mockRepo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, savedUser.ID)

	// Test error case
	mockRepo.mockDB.createError = errors.New("database error")
	err = mockRepo.Add(ctx, user)
	assert.Error(t, err)
}

func TestRepositoryExistsByEmail(t *testing.T) {
	// Setup
	mockRepo := NewMockRepository()
	ctx := context.Background()

	// Create test user
	user := createTestUser()
	err := mockRepo.Add(ctx, user)
	require.NoError(t, err)

	// Test exists
	exists, err := mockRepo.ExistsByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test does not exist
	exists, err = mockRepo.ExistsByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test error case
	mockRepo.mockDB.existsError = errors.New("database error")
	_, err = mockRepo.ExistsByEmail(ctx, user.Email)
	assert.Error(t, err)
}

func TestRepositoryGetByEmail(t *testing.T) {
	// Setup
	mockRepo := NewMockRepository()
	ctx := context.Background()

	// Create test user
	user := createTestUser()
	err := mockRepo.Add(ctx, user)
	require.NoError(t, err)

	// Test find by email
	foundUser, err := mockRepo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)

	// Test not found
	foundUser, err = mockRepo.GetByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, foundUser)

	// Test error case
	mockRepo.mockDB.readError = errors.New("database error")
	_, err = mockRepo.GetByEmail(ctx, user.Email)
	assert.Error(t, err)
}

func TestRepositoryGetByID(t *testing.T) {
	// Setup
	mockRepo := NewMockRepository()
	ctx := context.Background()

	// Create test user
	user := createTestUser()
	err := mockRepo.Add(ctx, user)
	require.NoError(t, err)

	// Test find by ID
	foundUser, err := mockRepo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)

	// Test not found
	foundUser, err = mockRepo.GetByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, foundUser)

	// Test error case
	mockRepo.mockDB.readError = errors.New("database error")
	_, err = mockRepo.GetByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestRepositoryUpdate(t *testing.T) {
	// Setup
	mockRepo := NewMockRepository()
	ctx := context.Background()

	// Create test user
	user := createTestUser()
	err := mockRepo.Add(ctx, user)
	require.NoError(t, err)

	// Update user fields
	user.FirstName = "Updated"
	user.LastName = "Name"

	// Test Update
	err = mockRepo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify user was updated
	updatedUser, err := mockRepo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updatedUser.FirstName)
	assert.Equal(t, "Name", updatedUser.LastName)

	// Test error case
	mockRepo.mockDB.updateError = errors.New("database error")
	err = mockRepo.Update(ctx, user)
	assert.Error(t, err)
}
