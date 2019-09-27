package main

import (
	"encoding/json"
)

type Packet struct {
	RPC      string
	IP       string
	NodeID   KademliaID
	contacts []Contact
}

func (packet *Packet) String() string {
	temp := ""
	for _, elem := range (*packet).contacts {
		temp += ", " + elem.String()
	}
	return "Header: \n" + "\tRPC: " + packet.RPC + "\n\tNodeID: " + packet.NodeID.String() + "\n\tIP: " + packet.IP + "\nPayload:\n\tContacts:" + temp
}

func EncodePacket(rpc string, nodeID KademliaID, ip string, contacts []Contact) []byte {
	packet := Packet{rpc, ip, nodeID, contacts}
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
