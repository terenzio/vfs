package errors

import (
	"fmt"
)

// GENERAL ERRORS ========================================

// ErrInvalidName ErrUserNotFound is an error that is returned when a user does not exist
func ErrInvalidName(name string) error {
	return fmt.Errorf("The name: %s contains invalid chars. Only alphabets and numbers are allowed.", name)
}

// ErrNameTooLong is an error that is returned when a username is too long
func ErrNameTooLong(name string) error {
	return fmt.Errorf("The name: %s is too long. The maximum length is 30 characters.", name)
}

// USER ERRORS ========================================

// ErrUserExists is an error that is returned when a user already exists
func ErrUserExists(username string) error {
	return fmt.Errorf("The user: %s already exists.", username)
}

// ErrUserNotFound is an error that is returned when a user does not exist
func ErrUserNotExists(username string) error {
	return fmt.Errorf("The user: %s doesn't exist.", username)
}

// FOLDER ERRORS ========================================

// ErrFolderExists is an error that is returned when a folder already exists
func ErrFolderExists(folderName string) error {
	return fmt.Errorf("The folder: %s already exists.", folderName)
}

// ErrFolderNotFound is an error that is returned when a folder does not exist
func ErrFolderNotFound(folderName string) error {
	return fmt.Errorf("The folder: %s doesn't exist.", folderName)
}
