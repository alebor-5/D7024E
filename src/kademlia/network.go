package main

import (
	"net"
	"strings"
	"fmt"
	"log"
)

type Network struct {
}

func handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	fmt.Println("Received from ", addr)
	fmt.Println("Received the following packet :")
	recPacket := DecodePacket(buffer[:n])
	fmt.Println(recPacket.String())
	if err != nil {
		log.Fatal(err)
	}

}

func Listen(ip string, port string) {
	// TODO
	udpAddr, err := net.ResolveUDPAddr("udp4", ip + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on " + ip + ":" + port)

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
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
	ip := "0.0.0.0"
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
