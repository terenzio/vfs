// domain/user.go

package models

// User represents the user entity in the domain layer.
type User struct {
	Username string
}

// UserRepository is the interface that wraps the basic user repository operations.
type UserRepository interface {
	Register(user User) error
	Exists(username string) (bool, error)
	ValidateUsername(username string) error
}
