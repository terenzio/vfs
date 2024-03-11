// domain/user.go

package models

// User represents the user entity in the domain layer.
type User struct {
	Username string
}

// UserRepository is the interface that wraps the basic user repository operations.
// It is used to abstract the storage of user data.
// The methods in this interface can be implemented by different storage types, like a database or an in-memory store.
// The UserRepository interface is defined in the domain layer, and the actual implementation is in the infrastructure layer.
// This separation allows the domain layer to be independent of the storage mechanism.
// The UserRepository interface is used by the application layer to interact with the user data.
// The application layer does not need to know the details of the storage mechanism, as it only interacts with the UserRepository interface.
// This allows for easy swapping of storage mechanisms without affecting the application layer.
type UserRepository interface {
	Register(user User) error
	Exists(username string) (bool, error)
	ValidateUsername(username string) error
}
