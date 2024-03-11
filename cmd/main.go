package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	timeString := now.Format(time.DateTime)

	fmt.Println("Hello IsCoolLab! Lets build a Virtual File System!")
	fmt.Println("Current time in string format:", timeString)
}
