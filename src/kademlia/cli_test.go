package main

import (
	"bytes"
	"fmt"
)

func ExampleCLI() {
	fmt.Println("ExampleCLI")
	node := InitKademliaNode()
	network := Network{&node}
	var stdin bytes.Buffer

	inputs := [7]string{
		"put \"hash\"\n", "put\n",
		"get \"10.10.1.3:4000\"", "get\n",
		"exit\n", "exit now\n",
		"superrandomcommandthatdoesnotexist"}
	for _, str := range inputs {
		stdin.Write([]byte(str))
		cli(&stdin, network)
	}

}
