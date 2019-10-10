package main

import (
	"testing"
)

func TestJoinNetwork(t *testing.T) {
	JoinNetwork(NewKademliaID("630f496549240d231d61365161424d442c040764"), "10.0.1.2")
	//InitKademliaNode()
}
