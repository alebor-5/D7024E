package main

import (
	"fmt"
	"time"
	"log"
	"net"
)
// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func handleUDPConnection(conn *net.UDPConn) {

         // here is where you want to do stuff like read or write to client

         buffer := make([]byte, 1024)

         n, addr, err := conn.ReadFromUDP(buffer)

         fmt.Println("UDP client : ", addr)
         fmt.Println("Received from UDP client :  ", string(buffer[:n]))

         if err != nil {
                 log.Fatal(err)
         }

         // NOTE : Need to specify client address in WriteToUDP() function
         //        otherwise, you will get this error message
         //        write udp : write: destination address required if you use Write() function instead of WriteToUDP()

         // write message back to client
         message := []byte("Hello UDP client!")
         _, err = conn.WriteToUDP(message, addr)

         if err != nil {
                 log.Println(err)
         }

 }

func server() {
	outboundIP := GetOutboundIP()
	fmt.Println(outboundIP)
	 hostName := "localhost" 
         portNum := "6000"
         service := hostName + ":" + portNum

         udpAddr, err := net.ResolveUDPAddr("udp4", service)

         if err != nil {
                 log.Fatal(err)
         }

         // setup listener for incoming UDP connection
         ln, err := net.ListenUDP("udp", udpAddr)

         if err != nil {
                 log.Fatal(err)
         }

         fmt.Println("UDP server up and listening on port 6000")
	 fmt.Println("The IP is " + hostName)

         defer ln.Close()

         for {
                 // wait for UDP client to connect
                 handleUDPConnection(ln)
         }
}


func client() {
	 fmt.Println("ENter IP")
	 var hostName string
	 fmt.Scanln(&hostName)

         portNum := "6000"

         service := hostName + ":" + portNum

         RemoteAddr, err := net.ResolveUDPAddr("udp", service)

         //LocalAddr := nil
         // see https://golang.org/pkg/net/#DialUDP

         conn, err := net.DialUDP("udp", nil, RemoteAddr)

         // note : you can use net.ResolveUDPAddr for LocalAddr as well
         //        for this tutorial simplicity sake, we will just use nil

         if err != nil {
                 log.Fatal(err)
         }

         log.Printf("Established connection to %s \n", service)
         log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
         log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

         defer conn.Close()

         // write a message to server
         message := []byte("Hello UDP server!")

         _, err = conn.Write(message)

         if err != nil {
                 log.Println(err)
         }

         // receive message from server
         buffer := make([]byte, 1024)
         n, addr, err := conn.ReadFromUDP(buffer)

         fmt.Println("UDP Server : ", addr)
         fmt.Println("Received from UDP server : ", string(buffer[:n]))

 }


func main(){
	for {
		fmt.Println("Hello, world!")
		time.Sleep(time.Second)
		fmt.Println("Enter 0 for client and 1 for server")
		var input string
		fmt.Scanln(&input)
		fmt.Println("You entered " + input)
		if input == "0" {
			client()
		}else{
			server()
		}
	}

}
