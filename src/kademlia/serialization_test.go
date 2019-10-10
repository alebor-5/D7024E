package main

import "testing"

func TestSerialization(t *testing.T) {
	rpc := "RPC"
	ID := *NewKademliaID("36551d16562d17404f37352c5857232c14119999")
	IP := "10.10.10.10"
	stringMsg := "Hello"
	contactsMsg := []Contact{NewContact(NewKademliaID("36551d16562d17404f37352c5857232c14110000"), "7.7.7.7")}

	encodedPacket := EncodePacket(rpc, ID, IP, EncodeString(stringMsg))
	encodedContacts := EncodeContacts(contactsMsg)
	PacketToByteArr(Packet{rpc, IP, ID, EncodeString(stringMsg)})

	decodedPacket := DecodePacket(encodedPacket)
	decodedString := DecodeString(decodedPacket.Message)
	decodedContacts := DecodeContacts(encodedContacts)

	if decodedPacket.RPC != rpc {
		t.Error("Expected: " + rpc + " Got: " + decodedPacket.RPC)
	}
	if decodedPacket.NodeID.String() != ID.String() {
		t.Error("Expected: " + ID.String() + " Got: " + decodedPacket.NodeID.String())
	}
	if decodedPacket.IP != IP {
		t.Error("Expected: " + IP + " Got: " + decodedPacket.IP)
	}
	if decodedString != stringMsg {
		t.Error("Expected: " + stringMsg + " Got: " + decodedString)
	}
	if decodedContacts[0].String() != contactsMsg[0].String() {
		t.Error("Expected: " + contactsMsg[0].String() + " Got: " + decodedContacts[0].String())
	}
}
