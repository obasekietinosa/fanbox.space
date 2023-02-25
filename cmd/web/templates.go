package main

import (
	"html/template"
	"path/filepath"

	"www.fanbox.space/internal/models"
)

type templateData struct {
	Letter  *models.Letter
	Letters []*models.Letter
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.go.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./ui/html/base.go.html",
			"./ui/html/partials/nav.go.html",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
