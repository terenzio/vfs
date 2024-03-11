package service_test

import (
	"github.com/stretchr/testify/assert"
	customErrors "github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
	"github.com/terenzio/vfs/service"
	"testing"
)

// MockFolderRepository is updated to include ValidateFolderName
type MockFolderRepository struct {
	CreateFolderFunc       func(models.Folder) error
	DeleteFolderFunc       func(string, string) error
	RenameFolderFunc       func(string, string, string) error
	ListFoldersFunc        func(string, string, string) ([]models.Folder, error)
	ValidateFolderNameFunc func(string) error
}

func (m *MockFolderRepository) CreateFolder(folder models.Folder) error {
	return m.CreateFolderFunc(folder)
}

func (m *MockFolderRepository) DeleteFolder(username, folderName string) error {
	return m.DeleteFolderFunc(username, folderName)
}

func (m *MockFolderRepository) RenameFolder(username, folderName, newFolderName string) error {
	return m.RenameFolderFunc(username, folderName, newFolderName)
}

func (m *MockFolderRepository) ListFolders(username, sortField, sortOrder string) ([]models.Folder, error) {
	return m.ListFoldersFunc(username, sortField, sortOrder)
}

func (m *MockFolderRepository) ValidateFolderName(folderName string) error {
	return m.ValidateFolderNameFunc(folderName)
}

// TestCreateFolder tests the CreateFolder method of FolderService using table-driven tests
func TestCreateFolder(t *testing.T) {
	tests := []struct {
		name            string
		userName        string
		folderName      string
		description     string
		mockUserSetup   func(userRepo *MockUserRepository)
		mockFolderSetup func(folderRepo *MockFolderRepository)
		expectError     error
	}{
		{
			name:        "ValidFolderCreation",
			userName:    "testUser",
			folderName:  "testFolder",
			description: "A test folder",
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ValidateFolderNameFunc = func(string) error { return nil }
				folderRepo.CreateFolderFunc = func(models.Folder) error { return nil }
			},
			expectError: nil,
		},
		{
			name:        "UserDoesNotExist",
			userName:    "unknownUser",
			folderName:  "testFolder",
			description: "A test folder",
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return false, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ValidateFolderNameFunc = func(string) error { return nil }
				folderRepo.CreateFolderFunc = func(models.Folder) error { return nil }
			},
			expectError: customErrors.ErrUserNotExists("unknownUser"),
		},
		{
			name:        "InvalidFolderName",
			userName:    "testUser",
			folderName:  "invalid@folder",
			description: "A test folder",
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ValidateFolderNameFunc = func(string) error { return customErrors.ErrInvalidName("invalid@folder") }
				folderRepo.CreateFolderFunc = func(models.Folder) error { return nil }
			},
			expectError: customErrors.ErrInvalidName("invalid@folder"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFolderRepository := &MockFolderRepository{}
			tt.mockFolderSetup(mockFolderRepository)
			mockUserRepository := &MockUserRepository{}
			tt.mockUserSetup(mockUserRepository)
			folderService := service.NewFolderService(mockFolderRepository, mockUserRepository)

			err := folderService.CreateFolder(tt.userName, tt.folderName, tt.description)
			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFolderService(t *testing.T) {
	tests := []struct {
		name            string
		testFunc        func(t *testing.T, folderService *service.FolderService)
		mockUserSetup   func(userRepo *MockUserRepository)
		mockFolderSetup func(folderRepo *MockFolderRepository)
	}{
		{
			name: "DeleteExistingFolder",
			testFunc: func(t *testing.T, folderService *service.FolderService) {
				err := folderService.DeleteFolder("testUser", "testFolder")
				assert.NoError(t, err)
			},
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.DeleteFolderFunc = func(string, string) error { return nil }
			},
		},
		{
			name: "RenameExistingFolder",
			testFunc: func(t *testing.T, folderService *service.FolderService) {
				err := folderService.RenameFolder("testUser", "testFolder", "newTestFolder")
				assert.NoError(t, err)
			},
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.RenameFolderFunc = func(string, string, string) error { return nil }
			},
		},
		{
			name: "ListFolders",
			testFunc: func(t *testing.T, folderService *service.FolderService) {
				_, err := folderService.ListFolders("testUser", "", "")
				assert.NoError(t, err)
			},
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ListFoldersFunc = func(string, string, string) ([]models.Folder, error) { return nil, nil }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFolderRepository := &MockFolderRepository{}
			tt.mockFolderSetup(mockFolderRepository)
			mockUserRepository := &MockUserRepository{}
			tt.mockUserSetup(mockUserRepository)
			folderService := service.NewFolderService(mockFolderRepository, mockUserRepository)

			tt.testFunc(t, folderService)
		})
	}
}
