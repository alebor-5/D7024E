package main

import "os"

func main() {
	node := NewKademliaNode()
	node.PrintIP()

	cli(os.Stdin)
}
