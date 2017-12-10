package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"eveKillmailCrawler/staticData"
)

func NewWebServer() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/amarr", handleAmarr)
	http.HandleFunc("/minmatar", handleMinmatar)
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	nameSlice := make([]string, 0)

	for _, item := range database.GetMostLostShipSorted(0) {
		if staticData.GetCategoryIDFromTypeID(item.Id) == 6 {
			nameSlice = append(nameSlice, staticData.GetTypeIDName(item.Id))
		}
	}

	a, _ := json.Marshal(nameSlice)
	io.WriteString(w, string(a))
}

func handleAmarr(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "amarr rrr", r.URL.Path[1:])
}

func handleMinmatar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "minmatar rrr", r.URL.Path[1:])
}
