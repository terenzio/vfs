// repository/folder_repository.go

package repository

import (
	"encoding/json"
	"fmt"
	"github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type FileFolderRepository struct {
	filePath string
	mu       sync.Mutex // ensures thread-safe access to the file
}

type storedFolder struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewFileFolderRepository creates a new instance of FileFolderRepository
func NewFileFolderRepository(filePath string) *FileFolderRepository {
	return &FileFolderRepository{
		filePath: filePath,
	}
}

func (r *FileFolderRepository) loadFolders() ([]storedFolder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []storedFolder{}, nil // Return an empty slice if the file doesn't exist
	}

	data, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var folders []storedFolder
	if err := json.Unmarshal(data, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *FileFolderRepository) saveFolders(folders []storedFolder) error {
	data, err := json.Marshal(folders)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.filePath, data, 0644)
}

func (r *FileFolderRepository) CreateFolder(folder models.Folder) error {
	folders, err := r.loadFolders()
	if err != nil {
		return err
	}

	// Check for existing folder with same name under the same username
	for _, f := range folders {
		if strings.EqualFold(f.Name, folder.Name) && f.Username == folder.Username {
			return errors.ErrFolderExists(folder.Name)
		}
	}

	folders = append(folders, storedFolder{
		Name:        folder.Name,
		Description: folder.Description,
		Username:    folder.Username,
		CreatedAt:   folder.CreatedAt,
	})

	return r.saveFolders(folders)
}

func (r *FileFolderRepository) DeleteFolder(username, folderName string) error {
	folders, err := r.loadFolders()
	if err != nil {
		return err
	}

	found := false
	for i, f := range folders {
		if f.Username == username && strings.EqualFold(f.Name, folderName) {
			// Remove folder from slice
			folders = append(folders[:i], folders[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.ErrFolderNotFound(folderName)
	}

	return r.saveFolders(folders)
}

func (r *FileFolderRepository) RenameFolder(username, folderName, newFolderName string) error {
	folders, err := r.loadFolders()
	if err != nil {
		return err
	}

	for i, f := range folders {
		if f.Username == username && strings.EqualFold(f.Name, folderName) {
			// Check if new name already exists
			for _, f2 := range folders {
				if f2.Username == username && strings.EqualFold(f2.Name, newFolderName) {
					return errors.ErrFolderExists(newFolderName)
				}
			}

			folders[i].Name = newFolderName
			return r.saveFolders(folders)
		}
	}

	return errors.ErrFolderNotFound(folderName)
}

func (r *FileFolderRepository) ListFolders(username, sortField, sortOrder string) ([]models.Folder, error) {
	folders, err := r.loadFolders()
	if err != nil {
		return nil, err
	}

	// Filter folders by username
	var userFolders []models.Folder
	for _, f := range folders {
		if f.Username == username {
			userFolders = append(userFolders, models.Folder{
				Username:    f.Username,
				Name:        f.Name,
				Description: f.Description,
				CreatedAt:   f.CreatedAt,
			})
		}
	}

	if len(userFolders) == 0 {
		return nil, fmt.Errorf("no folders found for user %s", username)
	}

	// Sorting
	switch sortField {
	case "--sort-name":
		sort.Slice(userFolders, func(i, j int) bool {
			if sortOrder == "desc" {
				return userFolders[i].Name > userFolders[j].Name
			}
			return userFolders[i].Name < userFolders[j].Name
		})
	case "--sort-created":
		sort.Slice(userFolders, func(i, j int) bool {
			if sortOrder == "desc" {
				return userFolders[i].CreatedAt.After(userFolders[j].CreatedAt)
			}
			return userFolders[i].CreatedAt.Before(userFolders[j].CreatedAt)
		})
	default:
		// Default to sort by name in ascending order if no sort flag is provided
		sort.Slice(userFolders, func(i, j int) bool {
			return userFolders[i].Name < userFolders[j].Name
		})
	}

	return userFolders, nil
}

// ValidateFolderName checks if the folder name is valid.
// It must contain only alphabets (uppercase and lowercase) and numbers, no spaces.
func (r *FileFolderRepository) ValidateFolderName(folderName string) error {
	// Check the length of the folder name first
	if len(folderName) > 30 {
		return errors.ErrNameTooLong(folderName)
	}

	// Regular expression to match usernames containing only alphabets and numbers.
	validFolderNameRegex := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	if !validFolderNameRegex.MatchString(folderName) {
		return errors.ErrInvalidName(folderName)
	}

	return nil
}
