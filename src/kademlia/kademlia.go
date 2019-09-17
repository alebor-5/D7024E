package main

import (
	"fmt"
)

type Kademlia struct {
	id  KademliaID
	ip string
	routingTable RoutingTable
}

func InitKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := GetIP()
	rt := NewRoutingTable(NewContact(id,ip))
	return  Kademlia{*id, GetIP(), *rt }
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
func (kademlia *Kademlia) PrintIP(){
	fmt.Println(kademlia.ip)
}
