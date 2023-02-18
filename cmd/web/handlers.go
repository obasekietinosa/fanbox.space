package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func handleError(err error, w http.ResponseWriter) {
	log.Print(err.Error())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.go.html",
		"./ui/html/partials/nav.go.html",
		"./ui/html/pages/home.go.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		handleError(err, w)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		handleError(err, w)
		return
	}
}

func letterView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Viewing a specific letter with ID %d...", id)
}

func letterCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new letter..."))
}
