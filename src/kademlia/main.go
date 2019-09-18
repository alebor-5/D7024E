package main

import "os"

func main() {
	node := InitKademliaNode()
	node.PrintIP()
	go Listen()
	cli(os.Stdin)
}
