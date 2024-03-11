package errors

import (
	"fmt"
)

// USER ERRORS

// ErrUserExists is an error that is returned when a user already exists
func ErrUserExists(username string) error {
	return fmt.Errorf("The %s has already existed.", username)
}

// ErrInvalidUsername ErrUserNotFound is an error that is returned when a user does not exist
func ErrInvalidUsername(username string) error {
	return fmt.Errorf("The %s contains invalid chars. Only alphabets and numbers are allowed.", username)
}

// ErrUsernameTooLong is an error that is returned when a username is too long
func ErrUsernameTooLong(username string) error {
	return fmt.Errorf("The %s is too long. The maximum length is 50 characters.", username)
}
