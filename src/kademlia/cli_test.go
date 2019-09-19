package main

import (
	"bytes"
)

func Examplecli() {
	var stdin bytes.Buffer

	inputs := [7]string{
		"put \"hash\"\n", "put\n",
		"get \"10.10.1.3:4000\"", "get\n",
		"exit\n", "exit now\n",
		"superrandomcommandthatdoesnotexist"}
	for _, str := range inputs {
		stdin.Write([]byte(str))
		cli(&stdin)
	}

	// Output:
	// >Store isn't implemented :(
	// >put takes exactly 1 argument! e.g. [put "48656c6c6f2066726f6d20414"]
	// >LookupData isn't implemented :(
	// >get takes exactly 1 argument! e.g. [get "111.111.111:0000"]
	// >This will terminate the node. Continue? [y/n]: Termination aborted
	// >exit doesn't take any arguments!
	// >Unknown command: superrandomcommandthatdoesnotexist
}
