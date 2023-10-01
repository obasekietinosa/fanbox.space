package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"www.fanbox.space/internal/models"
	"www.fanbox.space/internal/validator"
	"www.fanbox.space/ui"
)

type composeData struct {
	From       string
	To         string
	Email      string
	Subject    string
	Content    string
	Salutation string
	validator.Validator
}

type templateData struct {
	CurrentYear int
	Letter      *models.Letter
	Letters     []*models.Letter
	Compose     composeData
}

func readableDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func truncateText(text string, maxLength int) string {
	if maxLength > len(text) {
		return text
	}
	return text[:strings.LastIndex(text[:maxLength], " ")] + "..."
}

var functions = template.FuncMap{
	"readableDate": readableDate,
	"truncateText": truncateText,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.go.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.go.html",
			"html/partials/*.go.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
