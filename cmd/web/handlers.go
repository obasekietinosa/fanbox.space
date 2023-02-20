package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"www.fanbox.space/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.go.html",
		"./ui/html/partials/nav.go.html",
		"./ui/html/pages/home.go.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	letters, err := app.letters.Latest(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, letter := range letters {
		app.infoLog.Printf("%+v\n", letter)
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) letterView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	letter, err := app.letters.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.go.html",
		"./ui/html/partials/nav.go.html",
		"./ui/html/pages/letters/view.go.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", letter)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) letterCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := "etinosa.obaseki@gmail.com"
		subject := "A test of our love"
		author := "Etin Obaseki"
		recipient := "Ebose Osolase"
		content := "Hello. \nI write this letter to inform you that I have been absolutely smitten by you"
		salutation := "Your lover"

		id, err := app.letters.Insert(email, subject, author, recipient, content, salutation)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/letters/view?id=%d", id), http.StatusSeeOther)

		return
	}
	w.Write([]byte("Create a new letter..."))
}
