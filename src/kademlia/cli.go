package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cli() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">")
		scanner.Scan()
		args := strings.Fields(scanner.Text())

		if len(args) > 0 {
			switch strings.ToLower(args[0]) {
			case "ping":
				fmt.Println("Pong!")
			case "join":
				fmt.Println("Welcome!")
			default:
				fmt.Println("Unknown command: " + args[0])
			}
		}
	}
}
