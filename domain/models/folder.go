package models

import (
	"time"
)

// Folder represents the folder entity in the domain layer
type Folder struct {
	Username    string
	Name        string
	Description string
	CreatedAt   time.Time
}

// FolderRepository is an interface that abstracts the methods for folder persistence
type FolderRepository interface {
	Exists(userName, folderName string) (bool, error)
	CreateFolder(folder Folder) error
	DeleteFolder(username, folderName string) error
	RenameFolder(username, folderName, newFolderName string) error
	ListFolders(username, sortField, sortOrder string) ([]Folder, error)
	ValidateFolderName(folderName string) error
}
