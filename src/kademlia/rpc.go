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
		fmt.Println("Got and ping")
		//network.kademlia.routingTable.mux.Lock()
		//network.AddToRoutingTable(contact)
		//network.kademlia.routingTable.mux.Unlock()
		return EncodePacket("PONG", network.kademlia.id, network.kademlia.ip, []Contact{}, "")
	case "FIND_NODE":
		fmt.Println(idstring)
		target := NewKademliaID(idstring)

		returnContacts := network.kademlia.routingTable.FindClosestContacts(target, K+1)

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
		network.AddToRoutingTable(contact)
	case "FIND_NODE_RESULT":
		network.AddToRoutingTable(contact)
	case "UNKNOWN":
		fmt.Println("The RPC you sent was a non-standard RPC. Please check that the RPC exists in the HandleRequest function.")
	default:
		fmt.Println("UNKNOWN RESPONSE RPC: " + packet.RPC)

	}
	return packet
}

// AddToRoutingTable follows the Kademlia Way to instert a contact into RT
func (network *Network) AddToRoutingTable(contact Contact) {
	bucketIndex := network.kademlia.routingTable.getBucketIndex(contact.ID)
	if network.kademlia.routingTable.buckets[bucketIndex].Len() >= bucketSize {
		leastRecentlySeen := network.kademlia.routingTable.buckets[bucketIndex].list.Back().Value.(Contact)
		response := network.sendUDP("PING", leastRecentlySeen.Address, []Contact{}, "")
		if response.RPC == "PONG" {
			network.kademlia.routingTable.mux.Lock()
			network.kademlia.routingTable.AddContact(leastRecentlySeen)
			network.kademlia.routingTable.mux.Unlock()
		} else {
			//Remove leastRecently
			network.kademlia.routingTable.mux.Lock()
			network.kademlia.routingTable.RemoveContact(leastRecentlySeen)
			network.kademlia.routingTable.AddContact(contact)
			network.kademlia.routingTable.mux.Unlock()
		}
	} else {
		network.kademlia.routingTable.mux.Lock()
		network.kademlia.routingTable.AddContact(contact)
		network.kademlia.routingTable.mux.Unlock()
	}

}
