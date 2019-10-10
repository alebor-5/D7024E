package main

import "testing"

func TestHandleRequest(t *testing.T) {
	node := InitKademliaNode()
	net := Network{&node}
	pacList := []Packet{
		Packet{"PING", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), []byte{}},
		Packet{"FIND_VALUE", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), EncodeString("630f496249240d231d61365161424d442c045555")},
		Packet{"FIND_NODE", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), EncodeContacts([]Contact{})},
		Packet{"STORE", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), []byte{20}},
		Packet{"superRandomRPC", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), []byte{0}},
	}
	for _, pac := range pacList {
		net.HandleRequest(pac)
	}
}

func TestHandleResponse(t *testing.T) {
	node := InitKademliaNode()
	net := Network{&node}
	pacList := []Packet{
		Packet{"PONG", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), []byte{}},
		Packet{"FIND_VALUE_RESULT_V", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), EncodeString("630f496249240d231d61365161424d442c045555")},
		Packet{"FIND_VALUE_RESULT_C", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), EncodeString("630f496249240d231d61365161424d442c045555")},
		Packet{"FIND_NODE_RESULT", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), EncodeContacts([]Contact{})},
		Packet{"UNKNOWN", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), []byte{20}},
		Packet{"superRandomRPC", "10.0.0.70", *NewKademliaID("630f496249240d231d61365161424d442c041111"), []byte{0}},
	}
	for _, pac := range pacList {
		net.HandleResponse(pac)
	}
}

func TestAddToRoutingTable(t *testing.T) {
	node := Kademlia{*NewKademliaID("630f496249240d231d61365161424d442c040761"), "10.10.10.10", *NewRoutingTable(NewContact(NewKademliaID("630f496249240d231d61365161424d442c040761"), "10.10.10.10")), NewDataStore()}
	net := Network{&node}
	conList := []Contact{
		NewContact(NewKademliaID("1d40453829504e44372e5d542d61632c36404d35"), "10.0.0.70"),
		NewContact(NewKademliaID("3e4035362154582e5136410f062f3e551f0a344f"), "10.0.0.70"),
		NewContact(NewKademliaID("2d5e4d525e26255b54032541114b1e0956283736"), "10.0.0.70"),
		NewContact(NewKademliaID("31070c583a51351a5021320702453f374b0d3427"), "10.0.0.70"),
		NewContact(NewKademliaID("061c45620c134a581b21196244481b2c60201f4d"), "10.0.0.70"),
		NewContact(NewKademliaID("2f1d5a0b32242b2b21421f5215513f003b410830"), "10.0.0.70"),
	}

	for _, con := range conList {
		net.AddToRoutingTable(con)
	}
}
