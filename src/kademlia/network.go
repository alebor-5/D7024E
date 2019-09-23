package main

import (
	"net"
	"strings"
	"fmt"
	"log"
)

type Network struct {
	kademlia *Kademlia
}


func (network *Network) sendUDP(method string, ip string, payload string){
	packet := EncodePacket(method,network.kademlia.id.String(),network.kademlia.ip,payload)
	RemoteAddr, err := net.ResolveUDPAddr("udp", ip + ":6000")
	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sending " + method + " to " + ip)

	defer conn.Close()

	_, err = conn.Write(packet)

	if err != nil {
		log.Println(err)
	}
	

	// This is to handle the response from 
	buffer := make([]byte, 1024)
	response, RemoteAddr, err := conn.ReadFromUDP(buffer)

	resPacket := DecodePacket(buffer[:response])
	network.HandleResponse(resPacket)
}



func (network *Network) handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	fmt.Println("Received from ", addr)
	fmt.Println("Received the following packet :")
	recPacket := DecodePacket(buffer[:n])
	fmt.Println(recPacket.String())
	if err != nil {
		log.Fatal(err)
	}
	network.HandleRequest(recPacket, conn, addr)
	//TODO Send response

}

func (network *Network) SendResponse(method string, payload string, conn *net.UDPConn, addr *net.UDPAddr){
	packet := EncodePacket(method,network.kademlia.id.String(),network.kademlia.ip,payload)
	_, _ = conn.WriteToUDP(packet, addr )
	
}

func (network *Network) Listen() {
	ip := GetIP()
	udpAddr, err := net.ResolveUDPAddr("udp4", ip + ":6000")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on " + ip + ":6000")

	defer ln.Close()

	for {
		network.handleUDPConnection(ln)
	}
}
func (network *Network) SendPingMessage(contact Contact) {
	// TODO
	go network.sendUDP("PING", contact.Address, "")
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
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
