package main

import (
	"fmt"
	"testing"
)

func TestJoinNetwork(t *testing.T) {
	fmt.Printf("*************TestJoinNetwork************")
	randomID := NewKademliaID("36551d16562d17404f37352c5857232c14114263")
	randomID1 := "46555d16562d17404f37352c5857232c14514263"
	network := JoinNetwork(randomID, "192.168.1.69")
	fmt.Println(network.kademlia.id.String())
	network.kademlia.LookupData(randomID1)

	network.SendPingMessage(NewContact(randomID, "0.0.0.0"))
	network.kademlia.StoreData([]byte{22})
	//struct tests:
	//InitKademliaNode()
}
