package main

import (
	"encoding/json"
)

type Packet struct {
	Header Header
	Payload Payload
}

type Header struct {
	RPC string
	NodeID string
	IP string
}

type Payload struct {
	Message string
}

func (header *Header) String() string {
  return "RPC:" + header.RPC + ", NodeID:" + header.NodeID + ", IP:" + header.IP
}

func (payload *Payload) String() string {
	return "Payload: " + payload.Message
}

func (packet *Packet) String() string {
	return "Header: \n" + "\tRPC: " + packet.Header.RPC + "\n\tNodeID: " + packet.Header.NodeID + "\n\tIP: " + packet.Header.IP + "\nPayload:\n\tMessage" + packet.Payload.Message
}

func SendPacket(rpc string, nodeID string, message string) []byte {
	header := Header{rpc, nodeID, "1.1.1.1"}
	payload := Payload{message}
	packet := Packet{header,payload}
	return packetToJSON(packet)
}

func DecodePacket(encodedPacket []byte) Packet {
	res := Packet{}
	json.Unmarshal(encodedPacket, &res)
	return res
}

func packetToJSON(packet Packet) []byte {
	js, _ := json.Marshal(packet)
	return js
}
