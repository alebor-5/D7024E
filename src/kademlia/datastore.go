package main

import "sync"

type DataStore struct {
	sync.Mutex
	m map[string][]byte
}

func NewDataStore() DataStore {
	ds := DataStore{}
	ds.m = make(map[string][]byte)
	return ds
}

func (ds DataStore) Insert(key string, val []byte) {
	ds.Lock()
	ds.m[key] = val
	ds.Unlock()
}

// GetIfExists returns the byte array belonging to the key if it exists,
// returns false otherwise.
func (ds DataStore) GetIfExists(key string) interface{} {
	ds.Lock()
	val, e := ds.m[key]
	ds.Unlock()
	if e {
		return val
	}
	return false
}
