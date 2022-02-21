package main

import (
	"html/template"
	"net/url"
	"path/filepath"
	"time"

	"joylanguageschool.ru/pkg/models"
)

type templateData struct {
	CSRFToken   string
	AuthenticatedUser *models.User
	CurrentYear int
	Flash       string
	FormData    url.Values
	FormErrors  map[string]string
	Post        *models.Post
	Posts       []*models.Post
}

// Display a readable human date
func humanDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Create a template cache for templates
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
