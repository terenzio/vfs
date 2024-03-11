package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	customErrors "github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
	"github.com/terenzio/vfs/service"
)

// MockFileRepository is a mock of FileRepository
type MockFileRepository struct {
	CreateFileFunc       func(models.File) error
	DeleteFileFunc       func(string, string, string) error
	ListFilesFunc        func(string, string, string, string) ([]models.File, error)
	ValidateFileNameFunc func(string) error
}

func (m *MockFileRepository) CreateFile(file models.File) error {
	return m.CreateFileFunc(file)
}

func (m *MockFileRepository) DeleteFile(userName, folderName, fileName string) error {
	return m.DeleteFileFunc(userName, folderName, fileName)
}

func (m *MockFileRepository) ListFiles(userName, folderName, sortField, sortOrder string) ([]models.File, error) {
	return m.ListFilesFunc(userName, folderName, sortField, sortOrder)
}

func (m *MockFileRepository) ValidateFileName(fileName string) error {
	return m.ValidateFileNameFunc(fileName)
}

// TestFileService_CreateFile tests the CreateFile method using table-driven tests
func TestCreateFile(t *testing.T) {
	tests := []struct {
		name            string
		userName        string
		folderName      string
		fileName        string
		description     string
		mockUserSetup   func(userRepo *MockUserRepository)
		mockFolderSetup func(folderRepo *MockFolderRepository)
		mockFileSetup   func(fileRepo *MockFileRepository)
		expectedError   error
	}{
		{
			name:        "ValidFileCreation",
			userName:    "testUser",
			folderName:  "testFolder",
			fileName:    "testFile",
			description: "A test file",
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ExistsFunc = func(string, string) (bool, error) { return true, nil }
			},
			mockFileSetup: func(fileRepo *MockFileRepository) {
				fileRepo.ValidateFileNameFunc = func(string) error { return nil }
				fileRepo.CreateFileFunc = func(models.File) error { return nil }
			},
			expectedError: nil,
		},
		{
			name:        "UserDoesNotExist",
			userName:    "unknownUser",
			folderName:  "testFolder",
			fileName:    "testFile",
			description: "A test file",
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return false, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ExistsFunc = func(string, string) (bool, error) { return true, nil }
			},
			mockFileSetup: func(fileRepo *MockFileRepository) {
				fileRepo.ValidateFileNameFunc = func(string) error { return nil }
				fileRepo.CreateFileFunc = func(models.File) error { return nil }
			},
			expectedError: customErrors.ErrUserNotExists("unknownUser"),
		},
		{
			name:        "FolderDoesNotExist",
			userName:    "testFolder",
			folderName:  "unknownFolder",
			fileName:    "testFile",
			description: "A test file",
			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ExistsFunc = func(string, string) (bool, error) { return false, nil }
			},
			mockFileSetup: func(fileRepo *MockFileRepository) {
				fileRepo.ValidateFileNameFunc = func(string) error { return nil }
				fileRepo.CreateFileFunc = func(models.File) error { return nil }
			},
			expectedError: customErrors.ErrFolderNotFound("unknownFolder"),
		},
		{
			name:        "InvalidFileName",
			userName:    "testUser",
			folderName:  "testFolder",
			fileName:    "invalid@file",
			description: "A test file",

			mockUserSetup: func(userRepo *MockUserRepository) {
				userRepo.ExistsFunc = func(string) (bool, error) { return true, nil }
			},
			mockFolderSetup: func(folderRepo *MockFolderRepository) {
				folderRepo.ExistsFunc = func(string, string) (bool, error) { return true, nil }
			},
			mockFileSetup: func(fileRepo *MockFileRepository) {
				fileRepo.ValidateFileNameFunc = func(string) error { return customErrors.ErrInvalidName("invalid@file") }
				fileRepo.CreateFileFunc = func(models.File) error { return nil }
			},
			expectedError: customErrors.ErrInvalidName("invalid@file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFolderRepository := &MockFolderRepository{}
			tt.mockFolderSetup(mockFolderRepository)
			mockUserRepository := &MockUserRepository{}
			tt.mockUserSetup(mockUserRepository)
			mockFileRepository := &MockFileRepository{}
			tt.mockFileSetup(mockFileRepository)
			fileService := service.NewFileService(mockFileRepository, mockFolderRepository, mockUserRepository)

			err := fileService.CreateFile(tt.userName, tt.folderName, tt.fileName, tt.description)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
