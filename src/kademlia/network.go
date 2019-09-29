package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Network struct {
	kademlia *Kademlia
}

func (network *Network) sendUDP(method string, ip string, contacts []Contact, message string) Packet {
	byteArr := EncodePacket(method, network.kademlia.id, network.kademlia.ip, contacts, message)
	RemoteAddr, _ := net.ResolveUDPAddr("udp", ip+":6000")
	conn, _ := net.DialUDP("udp", nil, RemoteAddr)
	conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer conn.Close()
	fmt.Println("Sending " + method)
	conn.Write(byteArr)

	// This is to handle the response from
	buffer := make([]byte, 1024)
	response, RemoteAddr, err := conn.ReadFromUDP(buffer)

	resPacket := DecodePacket(buffer[:response])
	if err != nil {
		resPacket = Packet{"TIMEOUT", ip, network.kademlia.id, []Contact{}, ""}
	}
	res := network.HandleResponse(resPacket)
	fmt.Println("Got a " + res.RPC + " response")
	return res
}

func (network *Network) handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, _ := conn.ReadFromUDP(buffer)
	recPacket := DecodePacket(buffer[:n])
	sendResponsePacket := network.HandleRequest(recPacket)
	//TODO Send response

	network.SendResponse(sendResponsePacket, conn, addr)

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
	fmt.Println("hej")
	go network.sendUDP("PING", contact.Address, []Contact{}, "")
}

func (network *Network) SendFindContactMessage(shortlist *Shortlist, c chan int, target *Contact) {
	var contact Contact
	var shortItemPtr *ShortlistItem

	(*shortlist).mux.Lock()
	for i, item := range (*shortlist).ls {
		if i >= K {
			return
		} else if !item.sent {
			(*shortlist).ls[i].sent = true
			shortItemPtr = &(*shortlist).ls[i]
			contact = item.contact
			break
		} else if i == len((*shortlist).ls)-1 && item.sent {
			c <- 0
			return
		}
	}
	(*shortlist).mux.Unlock()
	response := network.sendUDP("FIND_NODE", contact.Address, []Contact{}, (*target).ID.String())
	(*shortlist).mux.Lock()
	fmt.Println(response.RPC)
	if response.RPC == "UNKNOWN" || response.RPC == "TIMEOUT" {
		fmt.Println("NU TOG VI BORT EN JÄVEL!!!")
		(*shortlist).remove(contact)
	} else {
		(*shortItemPtr).visited = true

		for _, con := range response.Contacts {
			(*shortlist).insert(target, con)
		}
		// for _, elem := range (*shortlist).ls {
		// 	fmt.Println(elem.contact.String() + ", Visited:" + strconv.FormatBool(elem.visited) + ", sent:" + strconv.FormatBool(elem.sent))
		// }
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
