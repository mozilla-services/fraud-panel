package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func dfHeartbeat(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is fine"))
}

func dfLbHeartbeat(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is fine"))
}

func dfVersion(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("version.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load version.json: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
