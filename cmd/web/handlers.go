package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"www.fanbox.space/internal/models"
	"www.fanbox.space/internal/validator"
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
		From:       from,
		To:         to,
		Email:      email,
		Salutation: "Yours sincerely",
	}

	app.render(w, http.StatusOK, "compose.go.html", data)
}

func (app *application) letterCreatePost(w http.ResponseWriter, r *http.Request) {
	var form composeData
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Validator.CheckField(validator.Email(form.Email), "email", "Enter a valid email address")
	form.Validator.CheckField(validator.NotBlank(form.Subject), "subject", "Enter a subject for your letter")
	form.Validator.CheckField(validator.NotBlank(form.From), "from", "Your letter must have a sender")
	form.Validator.CheckField(validator.NotBlank(form.To), "to", "Your letter must have a recipient")
	form.Validator.CheckField(validator.NotBlank(form.Content), "content", "Enter the content of your letter")
	form.Validator.CheckField(validator.NotBlank(form.Salutation), "salutation", "Your letter must have a salutation")

	if !form.Validator.Valid() {
		data := app.newTemplateData(r)
		data.Compose = form
		app.render(w, http.StatusUnprocessableEntity, "compose.go.html", data)
		return
	}

	id, err := app.letters.Insert(form.Email, form.Subject, form.From, form.To, form.Content, form.Salutation)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/letters/view/%d", id), http.StatusSeeOther)
}
