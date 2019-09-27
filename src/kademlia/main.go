package main

import (
	"flag"
	"os"
)

func main() {
	isBootstrapPtr := flag.Bool("bootstrap", false, "If false the node will try to join the hardcoded bootstrap node at startup")
	bootstrapIP := "10.0.0.2"
	bootstrapID := NewKademliaID("630f496249240d231d61365161424d442c040761")
	var network Network

	flag.Parse()
	if *isBootstrapPtr {
		network = Network{&Kademlia{*bootstrapID, bootstrapIP, *NewRoutingTable(NewContact(bootstrapID, bootstrapIP))}}
		go network.Listen()
	} else {
		network = JoinNetwork(bootstrapID, bootstrapIP)
	}
	network.cliLoop(os.Stdin)
}
