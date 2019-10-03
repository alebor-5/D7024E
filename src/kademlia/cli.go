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
		case "getnodeid":
			fmt.Println("Your nodeID is: " + network.kademlia.id.String())
		case "ping":
			network.SendPingMessage(NewContact(NewRandomKademliaID(), args[1:len(args)-1]))
		case "find0":
			id := NewKademliaID("1111111100000000000000000000000000000000")
			closest := network.kademlia.LookupContact(id)
			for _, elem := range closest {
				fmt.Println(elem.String())
			}
		case "put":
			if strExp.MatchString(args) {
				hash := network.kademlia.StoreData([]byte(args[1 : len(args)-1]))
				fmt.Println("The data was stored at: " + hash)
			} else {
				fmt.Println("put takes exactly 1 argument! e.g. [put \"Save this string.\"]")
			}
		case "get":
			if strExp.MatchString(args) {
				hash := args[1 : len(args)-1]
				if len(hash) == 20 {
					res, gotVal := network.kademlia.LookupData(hash).([]byte)
					if gotVal {
						fmt.Println(string(res))
					} else {
						fmt.Println("Could not find the value belongin to the hash:\n" + hash)
					}
				} else {
					fmt.Println("The hash must be exactly 20 bytes long")
				}
				fmt.Println("LookupData isn't implemented :(")
			} else {
				fmt.Println("get takes exactly 1 argument! e.g. [get \"48656c6c6f2066726f6d20414\"]")
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
