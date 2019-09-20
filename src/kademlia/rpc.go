package main

import (
	"fmt"
)

type RPC struct {
	
}

func (network *Network) HandleResponse(packet Packet) {
	contact := NewContact(NewKademliaID(packet.Header.NodeID), packet.Header.IP)
	switch packet.Header.RPC{
		case "PING":
			network.handlePing(contact)
		case "PONG":
			network.handlePong(contact)
		default:
			fmt.Println("UNKNOWN RPC: " + packet.Header.RPC)
	}
}

func (network *Network) handlePing(contact Contact){
	network.kademlia.routingTable.AddContact(contact)
	network.SendPongMessage(contact.Address)
}

func (network *Network) handlePong(contact Contact){
	network.kademlia.routingTable.AddContact(contact)
}
