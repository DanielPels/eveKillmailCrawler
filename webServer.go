package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
)

func NewWebServer() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/amarr", handleAmarr)
	http.HandleFunc("/minmatar", handleMinmatar)
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	a, _ := json.Marshal(database.GetMostKillerShipSorted(0))
	io.WriteString(w, string(a))
}

func handleAmarr(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "amarr rrr", r.URL.Path[1:])
}

func handleMinmatar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "minmatar rrr", r.URL.Path[1:])
}
