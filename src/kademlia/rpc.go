package main

import (
	"fmt"
)

type RPC struct {
}

func (network *Network) HandleRequest(packet Packet) []byte {
	contact := NewContact(NewKademliaID(packet.NodeID.String()), packet.IP)
	switch packet.RPC {
	case "PING":
		network.kademlia.routingTable.AddContact(contact)
		return EncodePacket("PONG", network.kademlia.id, network.kademlia.ip, []Contact{})
	default:
		fmt.Println("UNKNOWN REQUEST RPC: " + packet.RPC + ", sending a default Message")
		return EncodePacket("UNKNOWN", network.kademlia.id, network.kademlia.ip, []Contact{})
	}
}

func (network *Network) HandleResponse(packet Packet) Packet {
	contact := NewContact(NewKademliaID(packet.NodeID.String()), packet.IP)
	switch packet.RPC {
	case "PONG":
		network.kademlia.routingTable.AddContact(contact)
	default:
		fmt.Println("UNKNOWN RESPONSE RPC: " + packet.RPC)
	}
	return packet
}
