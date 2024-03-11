package main

import (
	"bufio"
	"github.com/terenzio/vfs/repository"

	"fmt"
	"github.com/terenzio/vfs/service"
	"os"
	"strings"
	"time"
)

func main() {

	now := time.Now()
	timeString := now.Format(time.DateTime)

	// Create a new user service with a text file repository for back up
	// Can be used for future implementation of a database repository
	userRepo := repository.NewFileUserRepository("users.txt")
	userService := service.NewUserService(userRepo)

	// Create a new folder service with a text file repository for back up
	// Can be used for future implementation of a database repository
	folderRepo := repository.NewFileFolderRepository("folders.txt")
	folderService := service.NewFolderService(folderRepo, userRepo)

	// Create a new file service with a text file repository for back up
	// Can be used for future implementation of a database repository
	fileRepo := repository.NewFileRepository("files.txt")
	fileService := service.NewFileService(fileRepo, folderRepo, userRepo)

	fmt.Println("==== IsCoolLab: Virtual File System CLI ====")
	fmt.Println("The current time is:", timeString)
	fmt.Println("Type 'help' to see available commands.")

	// Create a new scanner to read user input from the command line
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("# ")

		if !scanner.Scan() {
			break // Exit the loop if an error occurs or EOF is reached
		}

		input := scanner.Text()
		args := strings.Fields(input)

		switch args[0] {
		case "exit":
			fmt.Println("Exiting program.")
			return

		// Display the available commands
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("register [username]")
			fmt.Println("create-folder [username] [foldername] [description]")
			fmt.Println("delete-folder [username] [foldername]")
			fmt.Println("rename-folder [username] [foldername] [new-folder-name]")
			fmt.Println("list-folders [username] [--sort-name|--sort-created] [asc|desc]")

		// User commands
		case "register":
			if len(args) != 2 {
				fmt.Println("Usage: register [username]")
				continue
			}
			username := args[1]
			err := userService.Register(username)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Add '%s' successfully.\n", username)
			}

		// Folder commands
		case "create-folder":
			// This is simplified; you'll need to adapt it based on your FolderService implementation
			if len(args) < 3 {
				fmt.Println("Usage: create-folder [username] [foldername] [description]")
				continue
			}
			username, folderName, description := args[1], args[2], strings.Join(args[3:], " ")
			err := folderService.CreateFolder(username, folderName, description)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Create '%s' successfully.\n", folderName)
			}
		case "delete-folder":
			if len(args) != 3 {
				fmt.Println("Usage: delete-folder [username] [foldername]")
				continue
			}
			err := folderService.DeleteFolder(args[1], args[2])
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Delete '%s' successfully.\n", args[2])
			}
		case "rename-folder":
			if len(args) != 4 {
				fmt.Println("Usage: rename-folder [username] [foldername] [new-folder-name]")
				continue
			}
			err := folderService.RenameFolder(args[1], args[2], args[3])
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Rename '%s' to '%s' successfully.\n", args[2], args[3])
			}
		case "list-folders":
			if len(args) < 2 {
				fmt.Println("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
				continue
			}
			sortField := ""
			sortOrder := "asc"
			if len(args) > 2 {
				sortField = args[2]
				if sortField != "--sort-name" && sortField != "--sort-created" {
					fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
					continue
				}
				if len(args) == 4 {
					sortOrder = args[3]
					if sortOrder != "asc" && sortOrder != "desc" {
						fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
						continue
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

		//File commands
		case "create-file":
			if len(args) < 4 {
				fmt.Println("Usage: create-file [username] [foldername] [filename] [description]")
				continue
			}
			description := ""
			if len(args) > 4 {
				description = strings.Join(args[4:], " ")
			}
			err := fileService.CreateFile(args[1], args[2], args[3], description)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Create '%s' in %s/%s successfully.\n", args[3], args[1], args[2])
			}

		case "delete-file":
			if len(args) != 4 {
				fmt.Println("Usage: delete-file [username] [foldername] [filename]")
				continue
			}
			err := fileService.DeleteFile(args[1], args[2], args[3])
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Delete '%s' in %s/%s successfully.\n", args[3], args[1], args[2])
			}
		case "list-files":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
				continue
			}
			username, folderName := args[1], args[2]
			sortField := ""
			sortOrder := "asc" // Default sorting order
			if len(args) > 3 {
				sortField = args[3]
				if sortField != "--sort-name" && sortField != "--sort-created" {
					fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
					continue
				}
				if len(args) == 5 {
					sortOrder = args[4]
					if sortOrder != "asc" && sortOrder != "desc" {
						fmt.Fprintln(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
						continue
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

		default:
			fmt.Println("Error: Unrecognized command. Type 'help' to see available commands.")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
		}

	}

}
