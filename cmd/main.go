package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/terenzio/vfs/repository"
	"github.com/terenzio/vfs/service"
)

func main() {
	userService, folderService, fileService := initializeServices()
	displayWelcomeMessage()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("# ")

		if !scanner.Scan() {
			handleExit()
			break // Exit the loop if an error occurs or EOF is reached
		}

		input := scanner.Text()
		if input == "exit" {
			handleExit()
			return
		}

		processCommand(input, userService, folderService, fileService)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
	}
}

// initializeServices creates new instances of the user, folder, and file services
func initializeServices() (*service.UserService, *service.FolderService, *service.FileService) {
	userRepo := repository.NewFileUserRepository("users.txt")
	folderRepo := repository.NewFileFolderRepository("folders.txt")
	fileRepo := repository.NewFileRepository("files.txt")

	userService := service.NewUserService(userRepo)
	folderService := service.NewFolderService(folderRepo, userRepo)
	fileService := service.NewFileService(fileRepo, folderRepo, userRepo)

	return userService, folderService, fileService
}

// displayWelcomeMessage prints a welcome message to the console
func displayWelcomeMessage() {
	fmt.Println("==== IsCoolLab: Virtual File System CLI ====")
	fmt.Println("The current time is:", time.Now().Format(time.DateTime))
	fmt.Println("Type 'help' to see available commands.")
}

// handleExit performs cleanup and exits the program
func handleExit() {
	filesToCleanup := []string{"users.txt", "folders.txt", "files.txt"}
	cleanup(filesToCleanup)
	fmt.Println("Removed all temp files.")
	fmt.Println("Exiting program.\nSee you next time!")
}

// cleanup removes the files specified in the input slice
func cleanup(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err == nil {
			fmt.Printf("Removing file %s ...\n", file)
		}
	}
}

// processCommand handles the user input and calls the appropriate service method
func processCommand(input string, userService *service.UserService, folderService *service.FolderService, fileService *service.FileService) {
	args := strings.Fields(input)

	switch args[0] {
	case "help":
		displayHelp()
	case "register":
		registerUser(args, userService)
	case "create-folder":
		createFolder(args, folderService)
	case "delete-folder":
		deleteFolder(args, folderService)
	case "rename-folder":
		renameFolder(args, folderService)
	case "list-folders":
		listFolders(args, folderService)
	case "create-file":
		createFile(args, fileService)
	case "delete-file":
		deleteFile(args, fileService)
	case "list-files":
		listFiles(args, fileService)
	default:
		fmt.Println("Error: Unrecognized command. Type 'help' to see available commands.")
	}
}

// displayHelp prints the available commands to the console
func displayHelp() {
	fmt.Println("Available commands:")
	fmt.Println("> register [username]")
	fmt.Println("> create-folder [username] [foldername] [description]?")
	fmt.Println("> delete-folder [username] [foldername]")
	fmt.Println("> list-folders [username] [--sort-name|--sort-created] [asc|desc]")
	fmt.Println("> rename-folder [username] [foldername] [new-folder-name]")
	fmt.Println("> create-file [username] [foldername] [filename] [description]?")
	fmt.Println("> delete-file [username] [foldername] [filename]")
	fmt.Println("> list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
	fmt.Println("> exit")
}

// registerUser registers a new user
func registerUser(args []string, userService *service.UserService) {
	if len(args) != 2 {
		fmt.Println("Usage: register [username]")
		return
	}
	username := args[1]
	err := userService.Register(username)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Add '%s' successfully.\n", username)
	}
}

// createFolder creates a new folder
func createFolder(args []string, folderService *service.FolderService) {
	if len(args) < 3 {
		fmt.Println("Usage: create-folder [username] [foldername] [description]?")
		return
	}
	username, folderName := args[1], args[2]
	description := ""
	if len(args) > 3 {
		description = strings.Join(args[3:], " ")
	}
	err := folderService.CreateFolder(username, folderName, description)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Create '%s' successfully.\n", folderName)
	}
}

// deleteFolder deletes an existing folder
func deleteFolder(args []string, folderService *service.FolderService) {
	if len(args) != 3 {
		fmt.Println("Usage: delete-folder [username] [foldername]")
		return
	}
	err := folderService.DeleteFolder(args[1], args[2])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Delete '%s' successfully.\n", args[2])
	}
}

// renameFolder renames an existing folder
func renameFolder(args []string, folderService *service.FolderService) {
	if len(args) != 4 {
		fmt.Println("Usage: rename-folder [username] [foldername] [new-folder-name]")
		return
	}
	err := folderService.RenameFolder(args[1], args[2], args[3])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Rename '%s' to '%s' successfully.\n", args[2], args[3])
	}
}

// listFolders lists all folders for a given user
func listFolders(args []string, folderService *service.FolderService) {
	if len(args) < 2 {
		fmt.Println("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
		return
	}
	sortField := ""
	sortOrder := "asc"
	if len(args) > 2 {
		sortField = args[2]
		if sortField != "--sort-name" && sortField != "--sort-created" {
			fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
			return
		}
		if len(args) == 4 {
			sortOrder = args[3]
			if sortOrder != "asc" && sortOrder != "desc" {
				fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
				return
			}
		}
	}
	// List the folders
	folders, err := folderService.ListFolders(args[1], sortField, sortOrder)
	if err != nil {

		// If no folders are found, print a warning
		if err.Error() == fmt.Sprintf("no folders found for user %s", args[1]) {
			fmt.Printf("Warning: The %s doesn't have any folders.\n", args[1])
		} else {
			fmt.Printf("Error: %s\n", err.Error())
		}

	} else {

		// Determine the maximum length of each field across all files
		maxFolderLen, maxDescLen, maxDateLen, maxUserLen := 0, 0, 0, 0
		for _, f := range folders {
			if len(f.Name) > maxFolderLen {
				maxFolderLen = len(f.Name)
			}
			if len(f.Description) > maxDescLen {
				maxDescLen = len(f.Description)
			}
			folderCreatedAt := f.CreatedAt.Format(time.DateTime)
			if len(folderCreatedAt) > maxDateLen {
				maxDateLen = len(folderCreatedAt)
			}
			if len(f.Username) > maxUserLen {
				maxUserLen = len(f.Username)
			}
		}

		// Print header
		headerFmt := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds\n", maxFolderLen, maxDescLen, maxDateLen, maxUserLen)
		fmt.Printf(headerFmt, "Name", "Description", "Created At", "User Name")
		fmt.Println(strings.Repeat("-", maxFolderLen+maxDescLen+maxDateLen+maxUserLen+20))

		for _, folder := range folders {
			fmt.Printf(headerFmt, folder.Name, folder.Description, folder.CreatedAt.Format(time.DateTime), folder.Username)
		}
	}
}

// createFile creates a new file
func createFile(args []string, fileService *service.FileService) {
	if len(args) < 4 {
		fmt.Println("Usage: create-file [username] [foldername] [filename] [description]?")
		return
	}
	username, folderName, fileName := args[1], args[2], args[3]
	description := ""
	if len(args) > 4 {
		description = strings.Join(args[4:], " ")
	}
	err := fileService.CreateFile(username, folderName, fileName, description)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Create '%s' in %s/%s successfully.\n", args[3], args[1], args[2])
	}
}

// deleteFile deletes an existing file
func deleteFile(args []string, fileService *service.FileService) {
	if len(args) != 4 {
		fmt.Println("Usage: delete-file [username] [foldername] [filename]")
		return
	}
	err := fileService.DeleteFile(args[1], args[2], args[3])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Delete '%s' in %s/%s successfully.\n", args[3], args[1], args[2])
	}
}

// listFiles lists all files for a given user and folder
func listFiles(args []string, fileService *service.FileService) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
		return
	}
	username, folderName := args[1], args[2]
	sortField := ""
	sortOrder := "asc" // Default sorting order
	if len(args) > 3 {
		sortField = args[3]
		if sortField != "--sort-name" && sortField != "--sort-created" {
			fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
			return
		}
		if len(args) == 5 {
			sortOrder = args[4]
			if sortOrder != "asc" && sortOrder != "desc" {
				fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
				return
			}
		}
	}
	// List the files
	files, err := fileService.ListFiles(username, folderName, sortField, sortOrder)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else if len(files) == 0 {
		fmt.Println("Warning: The folder is empty.")
	} else {

		// Determine the maximum length of each field across all files
		maxFileLen, maxFolderLen, maxDescLen, maxDateLen, maxUserLen := 0, 0, 0, 0, 0
		for _, f := range files {
			if len(f.Name) > maxFileLen {
				maxFileLen = len(f.Name)
			}
			if len(f.Description) > maxDescLen {
				maxDescLen = len(f.Description)
			}
			folderCreatedAt := f.CreatedAt.Format(time.DateTime)
			if len(folderCreatedAt) > maxDateLen {
				maxDateLen = len(folderCreatedAt)
			}
			if len(f.FolderName) > maxFolderLen {
				maxFolderLen = len(f.FolderName)
			}
			if len(f.Username) > maxUserLen {
				maxUserLen = len(f.Username)
			}
		}

		// Print header
		headerFmt := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", maxFileLen, maxDescLen, maxDateLen, maxFolderLen, maxUserLen)
		fmt.Printf(headerFmt, "Name", "Description", "Created At", "Folder", "User Name")
		fmt.Println(strings.Repeat("-", maxFolderLen+maxDescLen+maxDateLen+maxUserLen+20))

		for _, file := range files {
			fmt.Printf(headerFmt, file.Name, file.Description, file.CreatedAt.Format(time.DateTime), file.FolderName, file.Username)
		}

	}
}
