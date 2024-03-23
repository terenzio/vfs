package models

import (
	"time"
)

// In DDD, the domain layer contains the core business logic and models.
//Interfaces can be used to define the expected behaviors (services) of your domain entities,
//making the core logic agnostic to specific implementations.

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

// Interface Advantages:
// Loose Coupling: By relying on interfaces rather than concrete implementations, different layers of your application
//communicate through well-defined contracts, reducing dependencies between them.
// Ease of Testing: Interfaces make it straightforward to mock dependencies in unit tests,
//allowing you to test components in isolation.
// Flexibility: Changing or adding new behavior becomes easier since you only need to provide a new implementation of
//an interface without altering the clients that depend on it.
