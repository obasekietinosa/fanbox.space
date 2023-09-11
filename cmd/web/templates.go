package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"www.fanbox.space/internal/models"
	"www.fanbox.space/ui"
)

type composeData struct {
	From  string
	To    string
	Email string
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

var functions = template.FuncMap{
	"readableDate": readableDate,
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
