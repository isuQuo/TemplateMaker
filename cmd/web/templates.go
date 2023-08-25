package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/isuquo/templatemaker/internal/models"
	"github.com/isuquo/templatemaker/ui"
	"github.com/justinas/nosurf"
)

type templateData struct {
	Template        *models.Template
	Templates       map[string][]*models.Template
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	Assessment      []string
	Query           []string
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths with the
	// extension '.html'. This matches the naming pattern for our HTML page
	// templates.
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'index.html') from the full file path
		name := filepath.Base(page)

		patterns := []string{
			"html/layouts/main.html",
			"html/partials/*.html",
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

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}
