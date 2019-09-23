package main

import (
	"flag"
	"os"
)

func main() {
	isBootstrapPtr := flag.Bool("bootstrap", false, "If false the node will try to join the hardcoded bootstrap node at startup")
	bootstrapIP := "10.0.0.2"
	bootstrapID := NewKademliaID("630f496249240d231d61365161424d442c040761")
	bootstrapNode := Kademlia{*bootstrapID, bootstrapIP, *NewRoutingTable(NewContact(bootstrapID, bootstrapIP))}
	var network Network
	var node Kademlia

	flag.Parse()
	if *isBootstrapPtr {
		//TODO: Bootstrap setup
		node = bootstrapNode
		network = Network{&node}
		go network.Listen()
	} else {
		//TODO: Genral setup
		//TODO: Try to join bootstrap node at: 10.0.0.2
		node = InitKademliaNode()
		network = Network{&node}
		go network.Listen()
		network.SendPingMessage(NewContact(bootstrapID, bootstrapIP))
	}
	network.cliLoop(os.Stdin)
}
