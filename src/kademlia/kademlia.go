package main

import (
	"fmt"
)

type Kademlia struct {
	id KademliaID
	ip string
}

func NewKademliaNode() Kademlia {
	id := NewKademliaID("48656c6c6f2066726f6d2041444d466163746f72792e636f6d")
	return Kademlia{*id, GetIP()}
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
func (kademlia *Kademlia) PrintIP() {
	fmt.Println(kademlia.ip)
}
