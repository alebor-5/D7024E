package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleObjects(t *testing.T) {
	fmt.Println("*****REST_Test****")
	node := InitKademliaNode()
	net := Network{&node}
	w := httptest.NewRecorder()
	recList := [](*http.Request){
		httptest.NewRequest("GET", "/objects/", nil),
		httptest.NewRequest("GET", "/objects/gfhgf", nil),
		httptest.NewRequest("GET", "/objects/1111111100000000000000000000000000000000", nil),
		httptest.NewRequest("POST", "/objects", nil),
		httptest.NewRequest("POST", "/objects", bytes.NewReader([]byte{10, 11})),
		httptest.NewRequest("PUT", "/objects", nil),
	}

	for _, r := range recList {
		net.handleObjects(w, r)
	}

}
