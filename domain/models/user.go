// domain/user.go

package models

// User represents the user entity in the domain layer.
type User struct {
	Username string
}

// In DDD, the domain layer contains the core business logic and models.
//Interfaces can be used to define the expected behaviors (services) of your domain entities,
//making the core logic agnostic to specific implementations.

// UserRepository is the interface that wraps the basic user repository operations.
type UserRepository interface {
	Register(user User) error
	Exists(username string) (bool, error)
	ValidateUsername(username string) error
}

// Interface Advantages:
// Loose Coupling: By relying on interfaces rather than concrete implementations, different layers of your application
//communicate through well-defined contracts, reducing dependencies between them.
// Ease of Testing: Interfaces make it straightforward to mock dependencies in unit tests,
//allowing you to test components in isolation.
// Flexibility: Changing or adding new behavior becomes easier since you only need to provide a new implementation of
//an interface without altering the clients that depend on it.
