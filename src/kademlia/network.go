package main

import (
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

type Network struct {
	kademlia *Kademlia
}

func (network *Network) sendUDP(method string, ip string, message []byte) Packet {
	msgSize := int(math.Pow(2, math.Ceil(math.Log2(200+K*100))))
	byteArr := EncodePacket(method, network.kademlia.id, network.kademlia.ip, message)
	RemoteAddr, resAddErr := net.ResolveUDPAddr("udp", ip+":6000")
	conn, dialErr := net.DialUDP("udp", nil, RemoteAddr)
	conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer conn.Close()
	fmt.Println("Sending " + method + " to " + ip)
	conn.Write(byteArr)
	if resAddErr != nil {
		fmt.Println("Det blev en resAddErr: " + resAddErr.Error())
	}
	if dialErr != nil {
		fmt.Println("Det blev en dialErr: " + dialErr.Error())
	}

	// This is to handle the response from
	buffer := make([]byte, msgSize)
	response, RemoteAddr, err := conn.ReadFromUDP(buffer)
	resPacket := DecodePacket(buffer[:response])
	if err != nil {
		resPacket = Packet{"TIMEOUT", ip, network.kademlia.id, []byte{}}
	}
	res := network.HandleResponse(resPacket)
	fmt.Println("Got a " + res.RPC + " response from " + res.IP)
	return res
}

func (network *Network) handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, int(math.Pow(2, math.Ceil(math.Log2(200+K*100)))))
	n, addr, _ := conn.ReadFromUDP(buffer)
	recPacket := DecodePacket(buffer[:n])
	sendResponsePacket := network.HandleRequest(recPacket)
	//TODO Send response

	go network.SendResponse(sendResponsePacket, conn, addr)

}

func (network *Network) SendResponse(packet []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	_, _ = conn.WriteToUDP(packet, addr)

}

func (network *Network) Listen() {
	ip := GetIP()
	udpAddr, _ := net.ResolveUDPAddr("udp4", ip+":6000")
	ln, _ := net.ListenUDP("udp", udpAddr)
	defer ln.Close()
	for {
		network.handleUDPConnection(ln)
	}
}

func (network *Network) SendPingMessage(contact Contact) {
	message := EncodeString("Just a PING message")
	go network.sendUDP("PING", contact.Address, message)
}

func (network *Network) SendFindContactMessage(shortlist *Shortlist, c chan int, targetID *KademliaID) {
	var contact Contact
	var shortItemPtr *ShortlistItem

	(*shortlist).mux.Lock()
	for i, item := range (*shortlist).ls {
		if i >= K {
			(*shortlist).mux.Unlock()
			return
		} else if !item.sent {
			(*shortlist).ls[i].sent = true
			shortItemPtr = &(*shortlist).ls[i]
			contact = item.contact
			break
		} else if i == len((*shortlist).ls)-1 && item.sent {
			c <- 0
			(*shortlist).mux.Unlock()
			return
		}
	}
	(*shortlist).mux.Unlock()
	message := EncodeString((*targetID).String())

	response := network.sendUDP("FIND_NODE", contact.Address, message)
	(*shortlist).mux.Lock()
	if response.RPC == "UNKNOWN" || response.RPC == "TIMEOUT" {
		(*shortlist).remove(contact.ID)
	} else {
		(*shortItemPtr).visited = true
		receivedContacts := DecodeContacts(response.Message)
		for _, con := range receivedContacts {
			fmt.Println(con.String())
			(*shortlist).insert(targetID, con)
		}
	}
	(*shortlist).mux.Unlock()
	c <- 0
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func GetIP() string {
	ip := "localhost"
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		byNameInterface, _ := net.InterfaceByName(i.Name)
		if i.Name == "eth0" {
			addresses, _ := byNameInterface.Addrs()
			for _, v := range addresses {
				ip = strings.TrimSuffix(v.String(), "/24")
			}
		}

	}
	return ip
}
