package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Ghada-Emad1/SnippetBox/internal/models"
)

type templatesData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
	CurrentYear int
	Form any
	Flash string 
	IsAuthenticated bool
	CSRFToken string
}

func HumanDate(t time.Time) string{
	return t.Format("02 Jan 2006 at 15:04 PM")
}
var functions=template.FuncMap{
	"humanDate":HumanDate,
}
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		
		cache[name] = ts
	}
	return cache, nil
}
