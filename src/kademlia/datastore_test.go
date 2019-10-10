package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDS(t *testing.T) {
	fmt.Println("TestDS")
	ds := NewDataStore()
	key := "superkey"
	val := []byte{0, 0, 0}
	ds.Insert(key, val)
	res1, _ := ds.GetIfExists(key).([]byte)
	res2, _ := ds.GetIfExists("not a valid key").(bool)
	if !bytes.Equal(val, res1) {
		t.Error("Expected: " + string(val) + " Got: " + string(res1))
	}
	if res2 {
		t.Error("Expected: " + string(val) + " Got: " + string(res1))
	}
}
