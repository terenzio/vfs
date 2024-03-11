package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	now := time.Now()
	timeString := now.Format(time.DateTime)

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
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("register [username]")
		default:
			fmt.Println("Error: Unrecognized command. Type 'help' to see available commands.")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
		}

	}

}
