package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	fmt.Println("Start listen on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
