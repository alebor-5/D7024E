package main

import (
	"net"
	"strings"
)

type Network struct {
	kademlia *Kademlia
}

func (network *Network) sendUDP(method string, ip string, contacts []Contact) Packet {
	byteArr := EncodePacket(method, network.kademlia.id, network.kademlia.ip, contacts)
	RemoteAddr, _ := net.ResolveUDPAddr("udp", ip+":6000")
	conn, _ := net.DialUDP("udp", nil, RemoteAddr)
	defer conn.Close()
	_, _ = conn.Write(byteArr)

	// This is to handle the response from
	buffer := make([]byte, 1024)
	response, RemoteAddr, _ := conn.ReadFromUDP(buffer)

	resPacket := DecodePacket(buffer[:response])
	return network.HandleResponse(resPacket)
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
	go network.sendUDP("PING", contact.Address, []Contact{})
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
