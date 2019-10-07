package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
)

const alpha = 3
const K = 10

type Kademlia struct {
	id           KademliaID
	ip           string
	routingTable RoutingTable
}

type ShortlistItem struct {
	contact Contact
	sent    bool
	visited bool
}

type Shortlist struct {
	ls  []ShortlistItem
	mux sync.Mutex
}

func JoinNetwork(bootstrapID *KademliaID, bootstrapIP string) Network {
	node := InitKademliaNode()
	network := Network{&node}
	//node.routingTable.mux.Lock()
	//node.routingTable.AddContact(NewContact(bootstrapID, bootstrapIP))
	//node.routingTable.mux.Unlock()
	go network.Listen()
	//TODO: Run iterative FIND_NODE on self
	network.SendPingMessage(NewContact(bootstrapID, bootstrapIP))
	//TODO: Refresh all buckets further away than the closest neighbor
	node.Refresh()
	return network
}

func InitKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := GetIP()
	rt := NewRoutingTable(NewContact(id, ip))
	return Kademlia{*id, GetIP(), *rt}
}

func (kademlia *Kademlia) LookupContact(targetID *KademliaID) []Contact {
	net := Network{kademlia}
	shortlist := Shortlist{}
	initContacts := kademlia.routingTable.FindClosestContacts(targetID, alpha)
	c := make(chan int, alpha)

	for _, contact := range initContacts {
		shortlist.mux.Lock()
		shortlist.insert(targetID, contact)
		shortlist.mux.Unlock()
	}

	for i := 0; i < alpha; i++ {
		go net.SendFindContactMessage(&shortlist, c, targetID)
	}

	for !lookupDone(&shortlist) {
		<-c
		go net.SendFindContactMessage(&shortlist, c, targetID)
	}

	shortlist.mux.Lock()
	var length int
	if len(shortlist.ls) < K {
		length = len(shortlist.ls)
	} else {
		length = K
	}
	result := []Contact{}
	for i := 0; i < length; i++ {
		result = append(result, shortlist.ls[i].contact)
	}
	shortlist.mux.Unlock()

	return result
}

func lookupDone(shortlist *Shortlist) bool {
	shortlist.mux.Lock()
	for i, item := range (*shortlist).ls {
		if i >= K {
			break
		} else if !item.sent || !item.visited {
			shortlist.mux.Unlock()
			return false
		}
	}
	shortlist.mux.Unlock()
	return true
}

// Inserts item sorted by distance to target
func (shortlist *Shortlist) insert(target *KademliaID, contact Contact) {
	conDist := contact.ID.CalcDistance(target)
	for i, shortItem := range (*shortlist).ls {
		itemDist := shortItem.contact.ID.CalcDistance(target)
		if shortItem.contact.ID.Equals(contact.ID) {
			return
		} else if (*conDist).Less(itemDist) {
			fst := (*shortlist).ls[:i]
			lst := (*shortlist).ls[i:]
			mid := []ShortlistItem{ShortlistItem{contact, false, false}}
			(*shortlist).ls = append(fst, append(mid, lst...)...)
			return
		}
	}
	(*shortlist).ls = append((*shortlist).ls, ShortlistItem{contact, false, false})
}

func (shortlist *Shortlist) remove(id *KademliaID) {
	for i, shortItem := range (*shortlist).ls {
		if shortItem.contact.ID.Equals(id) {
			fst := (*shortlist).ls[:i]
			if i == len((*shortlist).ls)-1 {
				(*shortlist).ls = append(fst, (*shortlist).ls[i+1:]...)
			} else {
				(*shortlist).ls = fst
			}
			return
		}
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
func (kademlia *Kademlia) PrintIP() {
	fmt.Println(kademlia.ip)
}

func (kademlia *Kademlia) Refresh() {
	bitString := ""
	tempArr := [IDLength]string{}
	for i := 0; i <= IDLength-1; i++ {

		bitstring := strconv.FormatInt(int64(kademlia.id[i]), 2)
		for j := len(bitstring); j < 8; j++ {
			bitstring = "0" + bitstring

		}
		tempArr[i] = bitstring
		fmt.Println(bitString)
	}

	for i := 0; i < IDLength; i++ {
		temp := tempArr
		for j := 0; j < 8; j++ {
			if string(temp[i][j]) == "1" {
				temp[i] = temp[i][:j] + string("0") + temp[i][:j+1]
			} else {
				temp[i] = temp[i][:j] + string("1") + temp[i][:j+1]
			}
			hexArr := [IDLength]byte{}
			for y := 0; y < IDLength; y++ {
				t, _ := strconv.ParseUint(temp[y], 2, 64)
				b := byte(t)
				hexArr[y] = b
				//fmt.Println(t)
			}
			//fmt.Println(hexArr)
			bucketKademliaID := hex.EncodeToString(hexArr[0:IDLength])
			//fmt.Println(bucketKademliaID)
			go kademlia.LookupContact(NewKademliaID(bucketKademliaID))
		}

		//temp = string(temp)
		//temp1, _ := strconv.ParseUint(temp, 2, 64)
		//fmt.Println(temp1)
		//fmt.Println(hex.DecodeString(temp))
		//kadId := NewKademliaID(temp1)
		//fmt.Println(kadId.String())
		//kademlia.LookupContact()
	}
	fmt.Println(len(bitString))
}
