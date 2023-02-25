package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"www.fanbox.space/internal/models"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	letters, err := app.letters.Latest(35)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.go.html", &templateData{
		Letters: letters,
	})
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

	app.render(w, http.StatusOK, "view.go.html", &templateData{
		Letter: letter,
	})
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
