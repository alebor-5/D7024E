package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCLI(t *testing.T) {
	fmt.Println("*****CLI_Test****")
	node := InitKademliaNode()
	network := Network{&node}
	var stdin bytes.Buffer

	inputs := []string{
		"enablelog", "disablelog",
		"help", "ping \"10.0.0.2\"\n",
		"put \"hash\"\n", "put\n",
		"get \"1111111100000000000000000000000000000000\"", "get\n",
		"exit\n", "exit now\n",
		"superrandomcommandthatdoesnotexist"}
	for _, str := range inputs {
		stdin.Write([]byte(str))
		cli(&stdin, network)
	}
}
