package main

import (
	"encoding/json"
)

// Packet definition
// The structure of a packet to be send over a UDP connection
type Packet struct {
	RPC     string
	IP      string
	NodeID  KademliaID
	Message []byte
}

/*
func (packet *Packet) String() string {
	temp := ""
	for _, elem := range (*packet).Contacts {
		temp += ", " + elem.String()
	}
	return "Header: \n" + "\tRPC: " + packet.RPC + "\n\tNodeID: " + packet.NodeID.String() + "\n\tIP: " + packet.IP + "\nPayload:\n\tContacts:" + temp
}
*/

// EncodePacket returns a byte array
func EncodePacket(rpc string, nodeID KademliaID, ip string, message []byte) []byte {
	packet := Packet{rpc, ip, nodeID, message}
	udpPacket, _ := json.Marshal(packet)
	return udpPacket
}

// EncodeContacts returns a byte array
func EncodeContacts(contacts []Contact) []byte {
	udpContacts, _ := json.Marshal(contacts)
	return udpContacts
}

// EncodeString returns a byte array
func EncodeString(message string) []byte {
	udpMessage, _ := json.Marshal(message)
	return udpMessage
}

// PacketToByteArr simply returns a byte array from a packet
func PacketToByteArr(packet Packet) []byte {
	byteArr, _ := json.Marshal(packet)
	return byteArr
}

// DecodePacket simply decodes a byte array, however, the message still need to be decoded as well
func DecodePacket(encodedPacket []byte) Packet {
	res := Packet{}
	json.Unmarshal(encodedPacket, &res)
	return res
}

// DecodeContacts decodes a byte array into a list of contacts
func DecodeContacts(encodedContacts []byte) []Contact {
	res := []Contact{}
	json.Unmarshal(encodedContacts, &res)
	return res
}

// DecodeString decodes a byte array into a string
func DecodeString(encodedContacts []byte) string {
	res := ""
	json.Unmarshal(encodedContacts, &res)
	return res
}
