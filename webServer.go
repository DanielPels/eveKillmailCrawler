package main

import (
	"net/http"
	"fmt"
	"io"
)

func NewWebServer() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/amarr", handleAmarr)
	http.HandleFunc("/minmatar", handleMinmatar)
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "KAKA: %s!", r.URL.Path[1:])
	io.WriteString(w, "DFHJSDKJHFKDJSHFKD")
}

func handleAmarr(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "amarr rrr", r.URL.Path[1:])
}

func handleMinmatar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "minmatar rrr", r.URL.Path[1:])
}
