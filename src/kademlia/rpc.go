package main

import (
	"fmt"
	"strconv"
)

func (network *Network) HandleRequest(packet Packet) []byte {
	idstring := packet.Message
	//fmt.Println(idstring)
	contact := NewContact(&packet.NodeID, packet.IP)
	le := len(packet.Contacts)
	fmt.Println(packet.RPC + " : " + strconv.Itoa(le))
	switch packet.RPC {
	case "PING":
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
		return EncodePacket("PONG", network.kademlia.id, network.kademlia.ip, []Contact{}, "")
	case "FIND_NODE":
		fmt.Println(idstring)
		target := NewKademliaID(idstring)
		network.kademlia.routingTable.mux.Lock()
		returnContacts := network.kademlia.routingTable.FindClosestContacts(target, K+1)
		network.kademlia.routingTable.mux.Unlock()

		for i, elem := range returnContacts {
			if contact.ID.Equals(elem.ID) {
				fst := returnContacts[:i]
				if i != len(returnContacts)-1 {
					returnContacts = append(fst, returnContacts[i+1:]...)
				} else {
					returnContacts = fst
				}
				break
			}
		}
		if len(returnContacts) > K {
			returnContacts = returnContacts[:K]
		}

		return EncodePacket("FIND_NODE_RESULT", network.kademlia.id, network.kademlia.ip, returnContacts, "")
	default:
		fmt.Println("UNKNOWN REQUEST RPC: " + packet.RPC + ", sending a default Message")
		return EncodePacket("UNKNOWN", network.kademlia.id, network.kademlia.ip, []Contact{}, "")
	}
}

func (network *Network) HandleResponse(packet Packet) Packet {
	contact := NewContact(&packet.NodeID, packet.IP)
	switch packet.RPC {
	case "PONG":
		fmt.Println("Got a PONG")
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
	case "FIND_NODE_RESULT":
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
