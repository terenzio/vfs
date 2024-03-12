package service

import (
	"github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
)

// UserService handles the service logic for users
type UserService struct {
	repo models.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo models.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register registers a new user with the given username
// It returns an error if the username is invalid, already exists, or if the registration fails
func (s *UserService) Register(username string) error {
	if err := s.repo.ValidateUsername(username); err != nil {
		return err
	}

	exists, err := s.repo.Exists(username)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrUserExists(username)
	}

	user := models.User{Username: username}
	return s.repo.Register(user)
}
