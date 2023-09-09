package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/letters/view/", app.letterView)
	mux.HandleFunc("/letters/create/", app.letterCreate)

	standard := alice.New(app.recoverPanic, app.logger, secureHeaders)

	return standard.Then(mux)
}
