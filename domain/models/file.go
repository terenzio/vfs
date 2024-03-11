// domain/file.go

package models

import (
	"time"
)

// File represents a file in the VFS
type File struct {
	Username    string
	FolderName  string
	Name        string
	Description string
	CreatedAt   time.Time
}

// FileRepository is an interface that abstracts the methods for file persistence
type FileRepository interface {
	CreateFile(file File) error
	DeleteFile(username, folderName, fileName string) error
	ListFiles(username, folderName, sortField, sortOrder string) ([]File, error)
	ValidateFileName(folderName string) error
}
