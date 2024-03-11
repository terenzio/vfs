package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	customErrors "github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
	"github.com/terenzio/vfs/service"
)

// mockUserRepository provides a mock implementation of the models.UserRepository interface
type mockUserRepository struct {
	RegisterFunc         func(models.User) error
	ExistsFunc           func(string) (bool, error)
	ValidateUsernameFunc func(string) error
}

func (m *mockUserRepository) Register(user models.User) error {
	return m.RegisterFunc(user)
}

func (m *mockUserRepository) Exists(username string) (bool, error) {
	return m.ExistsFunc(username)
}

func (m *mockUserRepository) ValidateUsername(username string) error {
	return m.ValidateUsernameFunc(username)
}

// TestRegister tests the Register method of UserService using table-driven tests
func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		setupMock     func(repo *mockUserRepository)
		expectedError error
	}{
		{
			name:     "ErrorInvalidUsername",
			username: "invalid!!user",
			setupMock: func(repo *mockUserRepository) {
				repo.ValidateUsernameFunc = func(string) error { return customErrors.ErrInvalidUsername("invalid!!user") }
			},
			expectedError: customErrors.ErrInvalidUsername("invalid!!user"),
		},
		{
			name:     "ErrorUserExists",
			username: "existingUser",
			setupMock: func(repo *mockUserRepository) {
				repo.ValidateUsernameFunc = func(string) error { return nil }
				repo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			expectedError: customErrors.ErrUserExists("existingUser"),
		},
		{
			name:     "Success",
			username: "newUser",
			setupMock: func(repo *mockUserRepository) {
				repo.ValidateUsernameFunc = func(string) error { return nil }
				repo.ExistsFunc = func(string) (bool, error) { return false, nil }
				repo.RegisterFunc = func(models.User) error { return nil }
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.setupMock(mockRepo)
			userService := service.NewUserService(mockRepo)

			err := userService.Register(tt.username)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
