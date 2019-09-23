package main

import (
	"flag"
	"os"
)

func main() {
	isBootstrapPtr := flag.Bool("bootstrap", false, "If false the node will try to join the hardcoded bootstrap node at startup")
	// bootstrapIP := "10.0.0.2"
	// bootstrapID := NewKademliaID("630f496249240d231d61365161424d442c040761")
	// bootstrapNode := NewContact(bootstrapID, bootstrapIP)
	node := InitKademliaNode()
	node.PrintIP()
	network := Network{&node}
	go network.Listen()
	flag.Parse()
	if *isBootstrapPtr {
		//TODO: Bootstrap setup
	} else {
		//TODO: Genral setup
		//TODO: Try to join bootstrap node at: 10.0.0.2
		network.SendPingMessage("10.0.0.2")
	}
	network.cliLoop(os.Stdin)
}
