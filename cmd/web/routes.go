package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/letters/view/:id", app.letterView)

	router.HandlerFunc(http.MethodGet, "/compose", app.letterCreate)
	router.HandlerFunc(http.MethodPost, "/compose", app.letterCreatePost)

	standard := alice.New(app.recoverPanic, app.logger, secureHeaders)

	return standard.Then(router)
}
