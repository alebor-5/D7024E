package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
)

func (network Network) handleObjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if len(r.URL.Path) <= len("/objects/") {
			http.Error(w, "Missing key in path", http.StatusBadRequest)
		} else {
			hash := r.URL.Path[len("/objects/"):]
			decoded, err := hex.DecodeString(hash)
			if len(decoded) == 20 && err == nil {
				res, gotVal := network.kademlia.LookupData(hash).([]byte)
				if gotVal {
					w.Write(res)
				} else {
					http.Error(w, "No value belonging to the hash: "+hash, http.StatusNotFound)
				}
			} else {
				http.Error(w, "Invalid hash format: "+hash, http.StatusBadRequest)
			}
		}
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			if len(body) > 0 {
				hash := network.kademlia.StoreData(body)
				w.Header().Set("Location", "/objects/"+hash)
				w.WriteHeader(http.StatusCreated)
				w.Write(body)
			} else {
				http.Error(w, "Empty body", http.StatusBadRequest)
			}
		}
	default:
		http.Error(w, "Invalid method: "+r.Method, http.StatusBadRequest)
	}
}

func (network Network) httpListen() {
	http.HandleFunc("/objects/", network.handleObjects)
	http.HandleFunc("/objects", network.handleObjects)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
