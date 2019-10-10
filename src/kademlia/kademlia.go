package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"sync"
)

const alpha = 3
const K = 5

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
	//TODO: Run iterative FIND_NODE on self
	node.LookupContact(&node.id)
	//TODO: Refresh all buckets further away than the closest neighbor
	node.Refresh()
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
	res, isVal := kademlia.vs.GetIfExists(hash).([]byte)
	if isVal {
		return res
	}
	targetID := NewKademliaID(hash)
	net := Network{kademlia}
	shortlist := Shortlist{}
	initContacts := kademlia.routingTable.FindClosestContacts(targetID, alpha)
	c := make(chan interface{}, alpha)
	if len(initContacts) < 1 {
		return []Contact{}
	}
	for _, contact := range initContacts {
		shortlist.mux.Lock()
		shortlist.insert(targetID, contact)
		shortlist.mux.Unlock()
	}

	for i := 0; i < alpha; i++ {
		go net.SendFindDataMessage(&shortlist, c, targetID)
	}

	for !lookupDone(&shortlist) {
		res, isVal := (<-c).([]byte)
		shortlist.mux.Lock()
		Log("From LookupData: ")
		for _, item := range shortlist.ls {
			Log("IP: " + item.contact.Address + " Sent: " + strconv.FormatBool(item.sent) + " Visited: " + strconv.FormatBool(item.visited))
		}
		shortlist.mux.Unlock()
		if isVal {
			return res
		}
		go net.SendFindDataMessage(&shortlist, c, targetID)
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
	hashValue := hex.EncodeToString(sha1.New().Sum(data)[0:IDLength])
	targetID := NewKademliaID(hashValue)
	kademlia.routingTable.mux.Lock()
	contacts := kademlia.routingTable.FindClosestContacts(targetID, K)
	kademlia.routingTable.mux.Unlock()
	// TODO: contacts := kademlia.LookupContact(targetID)
	for _, c := range contacts {
		go net.SendStoreMessage(data, c)
	}
	return hashValue
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

// Refresh simply loops through the 160-bit version of the local KademliaID
// It then flips one bit each, starting from the LSB until it reaches the MSB.
// For each change, it does a iterative find_node on that KademliaID
// The emptyBucketFlag is a check so that a find_node isn't sent until a bucket is populated
func (kademlia *Kademlia) Refresh() {
	tempArr := [IDLength]string{}
	for i := 0; i <= IDLength-1; i++ {
		bitstring := strconv.FormatInt(int64(kademlia.id[i]), 2)
		for j := len(bitstring); j < 8; j++ {
			bitstring = "0" + bitstring

		}
		tempArr[i] = bitstring
	}
	//Here temp array is the binary list of the KademliaID
	emptyBucketFlag := true
	for i := IDLength - 1; i >= 0; i-- {
		for j := 7; j >= 0; j-- {
			temp := tempArr
			//This if statement simply flips the [i][j] bit
			if string(temp[i][j]) == "1" {
				temp[i] = temp[i][:j] + string("0") + temp[i][j+1:]
			} else {
				temp[i] = temp[i][:j] + string("1") + temp[i][j+1:]
			}
			hexArr := [IDLength]byte{}
			//This for-loop makes the binary array into the same format as the KademliaID
			for y := 0; y < IDLength; y++ {
				t, _ := strconv.ParseUint(temp[y], 2, 64)
				b := byte(t)
				hexArr[y] = b

			}

			bucketKademliaID := hex.EncodeToString(hexArr[0:IDLength])
			//This is to get the bucketindex for the new kademliaID:
			bucketIndex := kademlia.routingTable.getBucketIndex(NewKademliaID(bucketKademliaID))
			if kademlia.routingTable.buckets[bucketIndex].Len() != 0 || !emptyBucketFlag == true {
				emptyBucketFlag = false
				Log("Refresh: Now sending a PING to a KademliaID belonging to bucket " + strconv.Itoa(bucketIndex))
				go kademlia.LookupContact(NewKademliaID(bucketKademliaID))
			}
		}
	}
}
