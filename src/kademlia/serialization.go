package main

import (
	"fmt"
	"encoding/json"
)


type Header struct {
	RPC string
	NodeID string
	IP string
}

func (header *Header) String() string {
  return "RPC:" + header.RPC + ", NodeID:" + header.NodeID + ", IP:" + header.IP
}

func getJSON(header Header) []byte{
	js, _ := json.Marshal(header)
	fmt.Println(js)
	return js
}

func setJSON(de []byte) string {
	res := Header{}

	json.Unmarshal(de, &res)
	return res.String()
}

