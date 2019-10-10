package main

import (
	"container/list"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// HandleRequest handles the received packet depending on the RPS attribute of the Packet
func (network *Network) HandleRequest(packet Packet) []byte {
	idstring := packet.NodeID.String()
	contact := NewContact(&packet.NodeID, packet.IP)
	switch packet.RPC {
	case "PING":
		network.AddToRoutingTable(contact)
		fmt.Println("Got a ping")
		message := EncodeString("Placeholder Message")
		return EncodePacket("PONG", network.kademlia.id, network.kademlia.ip, message)
	case "FIND_VALUE":
		message := DecodeString(packet.Message)
		res, b := network.kademlia.vs.GetIfExists(message).([]byte)
		if b {
			return EncodePacket("FIND_VALUE_RESULT_V", network.kademlia.id, network.kademlia.ip, res)
		} else {
			//PUT THIS IN A FUNCTION
			target := NewKademliaID(idstring)
			network.kademlia.routingTable.mux.Lock()
			returnContacts := network.kademlia.routingTable.FindClosestContacts(target, K+1)
			network.kademlia.routingTable.AddContact(contact)
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
			message := EncodeContacts(returnContacts)
			return EncodePacket("FIND_VALUE_RESULT_C", network.kademlia.id, network.kademlia.ip, message)
		}
	case "FIND_NODE":
		network.AddToRoutingTable(contact)
		target := NewKademliaID(idstring)
		network.kademlia.routingTable.mux.Lock()
		returnContacts := network.kademlia.routingTable.FindClosestContacts(target, K+1)
		network.kademlia.routingTable.AddContact(contact)
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
		message := EncodeContacts(returnContacts)
		return EncodePacket("FIND_N0DE_RESULT", network.kademlia.id, network.kademlia.ip, message)
	case "STORE":
		message := packet.Message
		hashValue := hex.EncodeToString(sha1.New().Sum(message)[0:IDLength])
		network.kademlia.vs.Insert(hashValue, message)
		pongMessage := EncodeString("Pong Message for STORE RPC")
		return EncodePacket("PONG", network.kademlia.id, network.kademlia.ip, pongMessage)
	default:
		fmt.Println("UNKNOWN REQUEST RPC: " + packet.RPC + ", sending a default Message")
		message := EncodeString("I don't understand the following RPC:" + packet.RPC)
		return EncodePacket("UNKNOWN", network.kademlia.id, network.kademlia.ip, message)
	}
}

// HandleResponse is the method that is called whenever a sender receives a message from the original recipient
func (network *Network) HandleResponse(packet Packet) Packet {
	contact := NewContact(&packet.NodeID, packet.IP)
	switch packet.RPC {
	case "PONG":
		fmt.Println("Got a PONG")
		network.AddToRoutingTable(contact)
	case "FIND_N0DE_RESULT":
		network.AddToRoutingTable(contact)
	case "FIND_VALUE_RESULT_V":
		network.AddToRoutingTable(contact)
	case "FIND_VALUE_RESULT_C":
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
	fmt.Println("Will print now")
	bucketIndex := network.kademlia.routingTable.getBucketIndex(contact.ID)

	//Check if exists and if so, set it first. Done

	var element *list.Element
	for e := network.kademlia.routingTable.buckets[bucketIndex].list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID
		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}
	if network.kademlia.routingTable.buckets[bucketIndex].Len() >= bucketSize && element == nil {

		//fmt.Println(network.kademlia.routingTable.buckets[bucketIndex].Len())
		leastRecentlySeen := network.kademlia.routingTable.buckets[bucketIndex].list.Back().Value.(Contact)
		message := EncodeString("Ping message")
		response := network.sendUDP("PING", leastRecentlySeen.Address, message)
		if response.RPC != "PONG" {
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
