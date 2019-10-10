package main

import (
	"fmt"
	"testing"
)

func TestJoinNetwork(t *testing.T) {
	fmt.Printf("*************TestJoinNetwork************")
	JoinNetwork(NewRandomKademliaID(), "10.0.1.2")
	//InitKademliaNode()
}
