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
			case "put":
				if len(args) != 2 {
					fmt.Println("put takes exactly 1 argument! [put <file_content>]")
				} else {
					fmt.Println("Store isn't implemented :(")
					//TODO: Run Store, if successful print the object hash.
				}
			case "get":
				if len(args) != 2 {
					fmt.Println("get takes exactly 1 argument! [get <hash>]")
				} else {
					fmt.Println("LookupData isn't implemented :(")
					//TODO: Run LookupData, if successful print the object content.
				}
			case "exit":
				fmt.Print("This will terminate the node. Continue? [y/n]: ")
				scanner.Scan()
				if strings.ToLower(scanner.Text()) == "y" {
					fmt.Println("Exiting node...")
					os.Exit(0)
				}
				fmt.Println("Termination aborted")
			default:
				fmt.Println("Unknown command: " + args[0])
			}
		}
	}
}
