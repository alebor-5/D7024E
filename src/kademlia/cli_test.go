package main

import (
	"bytes"
)

func ExampleCLI() {
	node := InitKademliaNode()
	network := Network{&node}
	var stdin bytes.Buffer

	inputs := []string{
		"ip\n", "refresh\n",
		"getcontacts\n", "getnodeid",
		"ping \"10.0.0.2\"\n", "find0\n",
		"put \"hash\"\n", "put\n",
		"get \"1111111100000000000000000000000000000000\"", "get\n",
		"exit\n", "exit now\n",
		"superrandomcommandthatdoesnotexist"}
	for _, str := range inputs {
		stdin.Write([]byte(str))
		cli(&stdin, network)
	}

	// Output:
	// >Store isn't implemented :(
	// >put takes exactly 1 argument! e.g. [put "48656c6c6f2066726f6d20414"]
	// >LookupData isn't implemented :(
	// >get takes exactly 1 argument! e.g. [get "111.111.111:0000"]
	// >This will terminate the node. Continue? [y/N]: Termination aborted
	// >exit doesn't take any arguments!
	// >Unknown command: superrandomcommandthatdoesnotexist
}
