package main

import (
	"encoding/json"
)

type Packet struct {
	RPC      string
	IP       string
	NodeID   KademliaID
	Contacts []Contact
	Message  string
}

func (packet *Packet) String() string {
	temp := ""
	for _, elem := range (*packet).Contacts {
		temp += ", " + elem.String()
	}
	return "Header: \n" + "\tRPC: " + packet.RPC + "\n\tNodeID: " + packet.NodeID.String() + "\n\tIP: " + packet.IP + "\nPayload:\n\tContacts:" + temp
}

func EncodePacket(rpc string, nodeID KademliaID, ip string, contacts []Contact, message string) []byte {
	packet := Packet{rpc, ip, nodeID, contacts, message}
	udpPacket, _ := json.Marshal(packet)
	return udpPacket
}

func PacketToByteArr(packet Packet) []byte {
	byteArr, _ := json.Marshal(packet)
	return byteArr
}

func DecodePacket(encodedPacket []byte) Packet {
	res := Packet{}
	json.Unmarshal(encodedPacket, &res)
	return res
}
