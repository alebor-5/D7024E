package main

import (
	"flag"
	"os"
)

func main() {
	isBootstrapPtr := flag.Bool("bootstrap", false, "If false the node will try to join the hardcoded bootstrap node at startup")
	flag.Parse()
	if *isBootstrapPtr {
		//TODO: Bootstrap setup
	} else {
		//TODO: Genral setup
		//TODO: Try to join bootstrap node at: 10.0.0.2
	}
	node := InitKademliaNode()
	node.PrintIP()
	network := Network{&node}
	go network.Listen()
	network.cliLoop(os.Stdin)
}
