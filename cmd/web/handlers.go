package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"www.fanbox.space/internal/models"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	letters, err := app.letters.Latest(35)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Letters = letters

	app.render(w, http.StatusOK, "home.go.html", data)
}

func (app *application) letterView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(r)
	data.Letter = letter

	app.render(w, http.StatusOK, "view.go.html", data)
}

func (app *application) letterCreate(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	email := r.URL.Query().Get("email")

	if (len(from) == 0) || (len(to) == 0) || (len(email) == 0) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	data := app.newTemplateData(r)
	data.Compose = composeData{
		From:  from,
		To:    to,
		Email: email,
	}

	app.render(w, http.StatusOK, "compose.go.html", data)
}

func (app *application) letterCreatePost(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, fmt.Sprintf("/letters/view/%d", id), http.StatusSeeOther)
}
