package errors

import (
	"fmt"
)

// NAMING ERRORS ========================================

// ErrInvalidName ErrUserNotFound is an error that is returned when a user does not exist
func ErrInvalidName(name string) error {
	return fmt.Errorf("The name [%s] contains invalid chars. Only alphabets and numbers are allowed.", name)
}

// ErrNameTooLong is an error that is returned when a username is too long
func ErrNameTooLong(name string) error {
	return fmt.Errorf("The name [%s] is too long. The maximum length is 30 characters.", name)
}

// USER ERRORS ========================================

// ErrUserExists is an error that is returned when a user already exists
func ErrUserExists(username string) error {
	return fmt.Errorf("The user [%s] already exists.", username)
}

// ErrUserNotExists ErrUserNotFound is an error that is returned when a user does not exist
func ErrUserNotExists(username string) error {
	return fmt.Errorf("The user [%s] doesn't exist.", username)
}

// FOLDER ERRORS ========================================

// ErrFolderExists is an error that is returned when a folder already exists
func ErrFolderExists(folderName string) error {
	return fmt.Errorf("The folder [%s] already exists.", folderName)
}

// ErrFolderNotFound is an error that is returned when a folder does not exist
func ErrFolderNotFound(folderName string) error {
	return fmt.Errorf("The folder [%s] doesn't exist.", folderName)
}

// FILE ERRORS ========================================

// ErrFileExists is an error that is returned when a file already exists
func ErrFileExists(fileName string) error {
	return fmt.Errorf("The file [%s] already exists.", fileName)
}

// ErrFileNotFound is an error that is returned when a file does not exist
func ErrFileNotFound(fileName string) error {
	return fmt.Errorf("The file [%s] doesn't exist.", fileName)
}
