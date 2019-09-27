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
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
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
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
	case "UNKNOWN":
		fmt.Println("The RPC you sent was a non-standard RPC. Please check that the RPC exists in the HandleRequest function.")
	default:
		fmt.Println("UNKNOWN RESPONSE RPC: " + packet.RPC)

	}
	return packet
}
