package main

import (
	"bytes"
	"testing"
)

func TestCLI(t *testing.T) {
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
}
