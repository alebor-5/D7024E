package main

import (
	"fmt"
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
	return network
}

func InitKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := GetIP()
	rt := NewRoutingTable(NewContact(id, ip))
	return Kademlia{*id, GetIP(), *rt}
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	net := Network{kademlia}
	shortlist := Shortlist{}
	initContacts := kademlia.routingTable.FindClosestContacts((*target).ID, alpha)
	c := make(chan int, alpha)

	for _, contact := range initContacts {
		shortlist.mux.Lock()
		shortlist.insert(target, contact)
		shortlist.mux.Unlock()
	}

	for i := 0; i < alpha; i++ {
		go net.SendFindContactMessage(&shortlist, c, target)
	}

	for !lookupDone(&shortlist) {
		<-c
		go net.SendFindContactMessage(&shortlist, c, target)
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

	fmt.Println("\n\n\n" + "Här är resultatet:")
	for _, elem := range result {
		fmt.Println(elem.String() + ", Distance: " + elem.ID.CalcDistance(target.ID).String())
	}
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
func (shortlist *Shortlist) insert(target *Contact, contact Contact) {
	conDist := contact.ID.CalcDistance((*target).ID)
	for i, shortItem := range (*shortlist).ls {
		itemDist := shortItem.contact.ID.CalcDistance((*target).ID)
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

func (shortlist *Shortlist) remove(contact Contact) {
	for i, shortItem := range (*shortlist).ls {
		if shortItem.contact.ID.Equals(contact.ID) {
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
