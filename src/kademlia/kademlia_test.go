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
	network.kademlia.LookupData(randomID1)
	network.SendPingMessage(NewContact(randomID, "0.0.0.0"))
	network.kademlia.StoreData([]byte{22})
}

func TestShortList(t *testing.T) {
	targetID := NewKademliaID("36551d16562d17404f37352c5857232c14110000")
	C1 := NewContact(NewKademliaID("36551d16562d17404f37352c5857232c14119999"), "1.1.1.1")
	C2 := NewContact(NewKademliaID("36551d16562d17404f37352c5857232c14113333"), "2.2.2.2")
	shortlist := Shortlist{}
	shortlist.insert(targetID, C1)
	shortlist.insert(targetID, C2)
	shortlist.insert(targetID, C2)
}
