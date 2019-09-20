package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func (network *Network) cliLoop(stdin io.Reader) {
	for {
		cli(stdin, *network)
	}
}

func cli(stdin io.Reader, network Network) {
	reader := bufio.NewReader(stdin)
	cmdExp := regexp.MustCompile(`^\s*\S+($|\s)`)
	spExp := regexp.MustCompile(`^\s*$`)
	strExp := regexp.MustCompile(`^".*"$`)

	fmt.Print(">")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	if cmdExp.MatchString(input) {
		cmd := cmdExp.FindString(input)
		args := strings.TrimSpace(input[len(cmd):])
		cmd = strings.TrimSpace(cmd)

		switch strings.ToLower(cmd) {
			case "ip":
				fmt.Println("Your IP is: " + GetIP())
			case "getcontacts":
				network.kademlia.routingTable.getAllContacts()
			case "ping":
				if strExp.MatchString(args) {
					go network.SendPingMessage(args[1:len(args)-1])
				} else {
					fmt.Println("put takes exactly 1 argument! e.g. [ping \"192.168.1.69\"]")
				}
			case "put":
				if strExp.MatchString(args) {
					fmt.Println("Store isn't implemented :(")
					//fmt.Println("Value: " + args[1:len(args)-1])
					//TODO: Run Store, if successful print the object hash.
				} else {
					fmt.Println("put takes exactly 1 argument! e.g. [put \"48656c6c6f2066726f6d20414\"]")
				}
			case "get":
				if strExp.MatchString(args) {
					fmt.Println("LookupData isn't implemented :(")
					//fmt.Println("IP: " + args[1:len(args)-1])
					//TODO: Run LookupData, if successful print the object content.
				} else {
					fmt.Println("get takes exactly 1 argument! e.g. [get \"111.111.111:0000\"]")
				}
			case "exit":
				if spExp.MatchString(args) {
					fmt.Print("This will terminate the node. Continue? [Y/n]: ")
					answer, _ := reader.ReadString('\n')
					answer = strings.TrimSuffix(answer, "\n")
					if strings.ToLower(answer) == "y" {
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
