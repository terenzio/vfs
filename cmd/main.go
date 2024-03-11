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

		default:
			fmt.Println("Error: Unrecognized command. Type 'help' to see available commands.")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
		}

	}

}
