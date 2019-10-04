package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
)

const alpha = 3
const K = 10

type Kademlia struct {
	id           KademliaID
	ip           string
	routingTable RoutingTable
	vs           DataStore
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
	node.routingTable.mux.Lock()
	node.routingTable.AddContact(NewContact(bootstrapID, bootstrapIP))
	node.routingTable.mux.Unlock()
	go network.Listen()
	//Run iterative FIND_NODE on self
	contacts := node.LookupContact(&node.id)
	for _, c := range contacts {
		if !node.id.Equals(c.ID) {
			node.routingTable.mux.Lock()
			node.routingTable.AddContact(c)
			node.routingTable.mux.Unlock()
		}
	}
	//TODO???: Refresh all buckets further away than the closest neighbor
	return network
}

func InitKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := GetIP()
	rt := NewRoutingTable(NewContact(id, ip))
	return Kademlia{*id, GetIP(), *rt, NewDataStore()}
}

func (kademlia *Kademlia) LookupContact(targetID *KademliaID) []Contact {
	net := Network{kademlia}
	shortlist := Shortlist{}
	initContacts := kademlia.routingTable.FindClosestContacts(targetID, alpha)
	c := make(chan int, alpha)
	if len(initContacts) < 1 {
		return []Contact{}
	}
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

func (kademlia *Kademlia) LookupData(hash string) interface{} {
	targetID := NewKademliaID(hash)
	net := Network{kademlia}
	shortlist := Shortlist{}
	initContacts := kademlia.routingTable.FindClosestContacts(targetID, alpha)
	c := make(chan int, alpha)
	if len(initContacts) < 1 {
		return []Contact{}
	}
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

func (kademlia *Kademlia) StoreData(data []byte) string {
	net := Network{kademlia}
	hashValue := hex.EncodeToString(sha1.New().Sum(data))
	targetID := NewKademliaID(hashValue)
	contacts := kademlia.LookupContact(targetID)
	for _, c := range contacts {
		go net.SendStoreMessage(data, c.ID)
	}
	return hashValue
}
func (kademlia *Kademlia) PrintIP() {
	fmt.Println(kademlia.ip)
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
