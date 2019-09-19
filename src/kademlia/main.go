package main

import "os"

func main() {
	node := NewKademliaNode()
	node.PrintIP()

	cliLoop(os.Stdin)
}
