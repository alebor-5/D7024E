package main

import "os"

func main() {
	//Todo make choice of either choosing to start or join network
	node := InitKademliaNode()
	node.PrintIP()
	network := Network{&node}
	go network.Listen()
	node.cli(os.Stdin)
}
