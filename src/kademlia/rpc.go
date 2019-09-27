package main

import (
	"fmt"
	"net"
)

type RPC struct {
}

func (network *Network) HandleRequest(packet Packet, conn *net.UDPConn, addr *net.UDPAddr) {
	contact := NewContact(NewKademliaID(packet.Header.NodeID), packet.Header.IP)
	switch packet.Header.RPC {
	case "PING":
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
		network.SendResponse("PONG", "", conn, addr)
	default:
		fmt.Println("UNKNOWN REQUEST RPC: " + packet.Header.RPC)
	}
}

func (network *Network) HandleResponse(packet Packet) {
	contact := NewContact(NewKademliaID(packet.Header.NodeID), packet.Header.IP)
	switch packet.Header.RPC {
	case "PONG":
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
	default:
		fmt.Println("UNKNOWN RESPONSE RPC: " + packet.Header.RPC)
	}
}
