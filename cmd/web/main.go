package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/letters/view", letterView)
	mux.HandleFunc("/letters/create", letterCreate)

	log.Print("Starting server on port 4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
