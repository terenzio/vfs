// repository/user_repository.go

package repository

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
)

// FileUserRepository handles the repository logic for users
type FileUserRepository struct {
	filePath string
	mu       sync.Mutex // ensures thread-safe access to the file
}

// NewFileUserRepository creates a new instance of a file-based user repository
func NewFileUserRepository(filePath string) *FileUserRepository {
	return &FileUserRepository{
		filePath: filePath,
	}
}

// Register adds a new user to the file
// The Register method adds a new user to the file. It takes a user model as an argument and writes the username to the file.
// The method uses a mutex to ensure that only one goroutine can write to the file at a time.
func (r *FileUserRepository) Register(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(user.Username + "\n")
	return err
}

// Exists checks if a username already exists in the file
// The Exists method checks if a username already exists in the file. It opens the file and scans through it line by line to find a match.
// It uses a mutex to ensure that only one goroutine can read from the file at a time.
func (r *FileUserRepository) Exists(username string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(r.filePath)
	if err != nil {
		// If the file doesn't exist, we treat it as no users exist yet.
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.EqualFold(scanner.Text(), username) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// ValidateUsername checks if the username is valid.
// It must contain only alphabets (uppercase and lowercase) and numbers, no spaces.
// The length of the username must be less than or equal to 30 characters.
func (r *FileUserRepository) ValidateUsername(username string) error {
	// Check the length of the username first
	if len(username) > 30 {
		return errors.ErrNameTooLong(username)
	}

	// Regular expression to match usernames containing only alphabets and numbers.
	validUsernameRegex := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	if !validUsernameRegex.MatchString(username) {
		return errors.ErrInvalidName(username)
	}

	return nil
}
