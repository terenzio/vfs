// repository/file_repository.go

package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"

	customErrors "github.com/terenzio/vfs/domain/errors"
	"github.com/terenzio/vfs/domain/models"
)

// FileRepository handles the repository logic for files
type FileRepository struct {
	filePath string
	mu       sync.Mutex // Ensures thread-safe access to the file
}

// storedFile represents the file structure stored in the file
type storedFile struct {
	Username    string `json:"username"`
	FolderName  string `json:"folderName"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

// NewFileRepository creates a new instance of FileRepository
func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{
		filePath: filePath,
	}
}

func (r *FileRepository) loadFiles() ([]storedFile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If the file does not exist, return an empty list
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []storedFile{}, nil
	}

	data, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var files []storedFile
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) saveFiles(files []storedFile) error {
	data, err := json.Marshal(files)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.filePath, data, 0644)
}

// CreateFile adds a new file to the repository
func (r *FileRepository) CreateFile(file models.File) error {
	files, err := r.loadFiles()
	if err != nil {
		return err
	}

	// Check if the file already exists within the same folder
	for _, f := range files {
		if f.Username == file.Username && f.FolderName == file.FolderName && f.Name == file.Name {
			return customErrors.ErrFileExists(file.Name)
		}
	}

	newFile := storedFile{
		Username:    file.Username,
		FolderName:  file.FolderName,
		Name:        file.Name,
		Description: file.Description,
		CreatedAt:   file.CreatedAt.Format("2006-01-02T15:04:05"),
	}

	files = append(files, newFile)

	return r.saveFiles(files)
}

// DeleteFile removes a file from the repository
func (r *FileRepository) DeleteFile(username, folderName, fileName string) error {
	files, err := r.loadFiles()
	if err != nil {
		return err
	}

	for i, f := range files {
		if f.Username == username && f.FolderName == folderName && f.Name == fileName {
			// Remove the file from the list
			files = append(files[:i], files[i+1:]...)
			return r.saveFiles(files)
		}
	}

	return customErrors.ErrFileNotFound(fileName)
}

// ListFiles returns a slice of files sorted based on the specified field and order.
func (r *FileRepository) ListFiles(username, folderName, sortField, sortOrder string) ([]models.File, error) {
	files, err := r.loadFiles()
	if err != nil {
		return nil, err
	}

	// Filter files by username and folderName
	var filteredFiles []storedFile
	for _, f := range files {
		if f.Username == username && f.FolderName == folderName {
			filteredFiles = append(filteredFiles, f)
		}
	}

	// Sorting
	switch sortField {
	case "--sort-name":
		sort.Slice(filteredFiles, func(i, j int) bool {
			if sortOrder == "desc" {
				return filteredFiles[i].Name > filteredFiles[j].Name
			}
			return filteredFiles[i].Name < filteredFiles[j].Name
		})
	case "--sort-created":
		sort.Slice(filteredFiles, func(i, j int) bool {
			if sortOrder == "desc" {
				return filteredFiles[i].CreatedAt > filteredFiles[j].CreatedAt
			}
			return filteredFiles[i].CreatedAt < filteredFiles[j].CreatedAt
		})
	default:
		// Default sorting by name in ascending order
		sort.Slice(filteredFiles, func(i, j int) bool {
			return filteredFiles[i].Name < filteredFiles[j].Name
		})
	}

	// Convert to domain.File slice
	var domainFiles []models.File
	for _, f := range filteredFiles {
		createdAt, err := time.Parse("2006-01-02T15:04:05", f.CreatedAt)
		if err != nil {
			// Handle the error, e.g., log it, skip this file, or use a zero time.
			// For this example, we'll log the error and continue with the next file.
			fmt.Printf("Error parsing date for file '%s': %v\n", f.Name, err)
			continue
		}

		domainFile := models.File{
			Username:    f.Username,
			FolderName:  f.FolderName,
			Name:        f.Name,
			Description: f.Description,
			CreatedAt:   createdAt,
		}
		domainFiles = append(domainFiles, domainFile)
	}

	return domainFiles, nil
}

// ValidateFileName checks if the folder name is valid.
// It must contain only alphabets (uppercase and lowercase) and numbers, no spaces.
// The length of the file name must be less than or equal to 30 characters.
func (r *FileRepository) ValidateFileName(fileName string) error {
	// Check the length of the folder name first
	if len(fileName) > 30 {
		return customErrors.ErrNameTooLong(fileName)
	}

	// Regular expression to match usernames containing only alphabets and numbers.
	validFolderNameRegex := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	if !validFolderNameRegex.MatchString(fileName) {
		return customErrors.ErrInvalidName(fileName)
	}

	return nil
}
