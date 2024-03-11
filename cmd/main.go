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

	fmt.Println("==== IsCoolLab: Virtual File System CLI ====")
	fmt.Println("The current time is:", timeString)
	fmt.Println("Type 'help' to see available commands.")

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
				if len(args) == 4 {
					sortOrder = args[3]
				}
			}
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

		default:
			fmt.Println("Error: Unrecognized command. Type 'help' to see available commands.")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
		}

	}

}
