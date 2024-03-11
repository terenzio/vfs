// service/file_service.go

package service

import (
	"time"

	"github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
)

// FileService handles the service logic for files
type FileService struct {
	fileRepo   models.FileRepository
	folderRepo models.FolderRepository
	userRepo   models.UserRepository
}

// NewFileService creates a new instance of FileService
func NewFileService(repo models.FileRepository, folderRepo models.FolderRepository, userRepo models.UserRepository) *FileService {
	return &FileService{fileRepo: repo, folderRepo: folderRepo, userRepo: userRepo}
}

// CreateFile creates a new file
func (s *FileService) CreateFile(userName, folderName, fileName, description string) error {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrUserNotExists(userName)
	}

	// Check if the folder exists
	exists, err = s.folderRepo.Exists(userName, folderName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrFolderNotFound(folderName)
	}

	// Check if the file fileName is valid
	if err := s.fileRepo.ValidateFileName(fileName); err != nil {
		return err
	}

	// Create the file
	file := models.File{
		Username:    userName,
		FolderName:  folderName,
		Name:        fileName,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return s.fileRepo.CreateFile(file)
}

// DeleteFile deletes a file
func (s *FileService) DeleteFile(userName, folderName, fileName string) error {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrUserNotExists(userName)
	}

	// Check if the folder exists
	exists, err = s.folderRepo.Exists(userName, folderName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrFolderNotFound(folderName)
	}

	// Delete the file
	return s.fileRepo.DeleteFile(userName, folderName, fileName)
}

// ListFiles lists the files in a folder
func (s *FileService) ListFiles(userName, folderName, sortField, sortOrder string) ([]models.File, error) {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.ErrUserNotExists(userName)
	}

	// Check if the folder exists
	exists, err = s.folderRepo.Exists(userName, folderName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.ErrFolderNotFound(folderName)
	}

	// List the files
	return s.fileRepo.ListFiles(userName, folderName, sortField, sortOrder)
}
