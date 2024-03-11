// service/folder_service.go

package service

import (
	"github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
	"time"
)

// FolderService handles the service logic for folders
type FolderService struct {
	folderRepo models.FolderRepository
	userRepo   models.UserRepository
}

// NewFolderService creates a new instance of FolderService
func NewFolderService(folderRepo models.FolderRepository, userRepo models.UserRepository) *FolderService {
	return &FolderService{folderRepo: folderRepo, userRepo: userRepo}
}

// CreateFolder creates a new folder
func (s *FolderService) CreateFolder(userName, folderName, description string) error {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrUserNotExists(userName)
	}

	// Check if the folder folderName is valid
	if err := s.folderRepo.ValidateFolderName(folderName); err != nil {
		return err
	}

	// Create the folder
	folder := models.Folder{
		Username:    userName,
		Name:        folderName,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return s.folderRepo.CreateFolder(folder)
}

// DeleteFolder deletes a folder
func (s *FolderService) DeleteFolder(userName, folderName string) error {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrUserNotExists(userName)
	}

	// Check if the folder folderName is valid
	if err := s.folderRepo.ValidateFolderName(folderName); err != nil {
		return err
	}

	// Delete the folder
	return s.folderRepo.DeleteFolder(userName, folderName)
}

// RenameFolder renames a folder
func (s *FolderService) RenameFolder(userName, folderName, newFolderName string) error {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrUserNotExists(userName)
	}

	// Rename the folder
	return s.folderRepo.RenameFolder(userName, folderName, newFolderName)
}

// ListFolders lists the folders
func (s *FolderService) ListFolders(userName, sortField, sortOrder string) ([]models.Folder, error) {

	// Check if the user exists
	exists, err := s.userRepo.Exists(userName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.ErrUserNotExists(userName)
	}

	// List the folders
	return s.folderRepo.ListFolders(userName, sortField, sortOrder)
}
