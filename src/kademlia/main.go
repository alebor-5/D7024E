package main

import (
	"time"
	"fmt"
)

func main() {
	for {
		time.Sleep(time.Second)
		/*fmt.Println("Enter 0 for client and 1 for server")
		var input string
		fmt.Scanln(&input)
		if input == "0" {
			client()
		}else{
			server()
		}*/
		node := InitKademliaNode()
		node.PrintIP()
		//c1 := NewContact(NewRandomKademliaID(),"localhost")
		//SendPONG(c1)	
		testPacket := DecodePacket(SendPacket("POST", "123456", "This is a message"))
		fmt.Println(testPacket.String())
	}

}
