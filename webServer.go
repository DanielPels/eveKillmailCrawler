package main

import (
	"net/http"
	"fmt"
	"eveKillmailCrawler/staticData"
	"eveKillmailCrawler/market"
	"encoding/json"
	"io"
)

func NewWebServer() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/amarr/", handleAmarr)
	http.HandleFunc("/minmatar", handleMinmatar)
	http.ListenAndServe(":8080", nil)
}

type tempstruct struct {
	Name  string  `json:"Name"`
	Price float64 `json:"Price"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	var nameSlice []tempstruct

	for _, item := range database.GetMostLostShipSorted(0) {
		if staticData.GetCategoryIDFromTypeID(item.Id) == 6 {
			nameSlice = append(nameSlice, tempstruct{
				Name:  staticData.GetTypeIDName(item.Id),
				Price: market.GetPriceOfTypeID(item.Id),
			})
		}
	}

	a, _ := json.Marshal(nameSlice)
	io.WriteString(w, string(a))
}

func handleAmarr(w http.ResponseWriter, r *http.Request) {
	fmt.Println()

	fmt.Fprintf(w, "amarr rrr", r.URL.Path[1:])
}

func handleMinmatar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "minmatar rrr", r.URL.Path[1:])
}
