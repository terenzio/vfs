// domain/file.go

package models

import (
	"time"
)

type File struct {
	Username    string
	FolderName  string
	Name        string
	Description string
	CreatedAt   time.Time
}

type FileRepository interface {
	CreateFile(file File) error
	DeleteFile(username, folderName, fileName string) error
	ListFiles(username, folderName, sortField, sortOrder string) ([]File, error)
	ValidateFileName(folderName string) error
}
