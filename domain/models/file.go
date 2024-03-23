// domain/file.go

package models

import (
	"time"
)

// In DDD, the domain layer contains the core business logic and models.
//Interfaces can be used to define the expected behaviors (services) of your domain entities,
//making the core logic agnostic to specific implementations.

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

// Interface Advantages:
// Loose Coupling: By relying on interfaces rather than concrete implementations, different layers of your application
//communicate through well-defined contracts, reducing dependencies between them.
// Ease of Testing: Interfaces make it straightforward to mock dependencies in unit tests,
//allowing you to test components in isolation.
// Flexibility: Changing or adding new behavior becomes easier since you only need to provide a new implementation of
//an interface without altering the clients that depend on it.
