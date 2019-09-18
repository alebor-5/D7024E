package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func cli() {
	scanner := bufio.NewScanner(os.Stdin)
	cmdExp := regexp.MustCompile(`^\s*\S+($|\s)`)
	spExp := regexp.MustCompile(`^\s*$`)
	strExp := regexp.MustCompile(`^".*"$`)
	for {
		fmt.Print(">")
		scanner.Scan()
		input := scanner.Text()

		if cmdExp.MatchString(input) {
			rawcmd := cmdExp.FindString(input)
			args := strings.TrimSpace(input[len(rawcmd):])
			cmd := strings.TrimSpace(rawcmd)

			switch strings.ToLower(cmd) {
			case "put":
				if strExp.MatchString(args) {
					fmt.Println("Store isn't implemented :(")
					fmt.Println("IP: " + args[1:len(args)-1])
					//TODO: Run Store, if successful print the object hash.
				} else {
					fmt.Println("put takes exactly 1 argument! e.g. [put \"111.111.111:0000\"]")
				}
			case "get":
				if strExp.MatchString(args) {
					fmt.Println("LookupData isn't implemented :(")
					fmt.Println("Key: " + args[1:len(args)-1])
					//TODO: Run LookupData, if successful print the object content.
				} else {
					fmt.Println("get takes exactly 1 argument! e.g. [get \"48656c6c6f2066726f6d20414\"]")
				}
			case "exit":
				if spExp.MatchString(args) {
					fmt.Print("This will terminate the node. Continue? [y/n]: ")
					scanner.Scan()
					if strings.ToLower(scanner.Text()) == "y" {
						fmt.Println("Exiting node...")
						os.Exit(0)
					}
					fmt.Println("Termination aborted")
				} else {
					fmt.Println("exit doesn't take any arguments!")
				}
			default:
				fmt.Println("Unknown command: " + cmd)
			}

		}
	}
}
