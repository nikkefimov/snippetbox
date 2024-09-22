package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox/pkg/models"
	"snippetbox/ui"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool   // Add a isAuthenticated field to the templateData struct.
	CSRFToken       string // Add a CSRFToken field.
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	// Conber the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template, FuncMap object and store it in a global variable
// this is essentially a string-keyd map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Initialize a map which keeps cache.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to het a slice of all filepaths in the ui.Fles embedded
	// filesystem which match the pattern 'html/pages/*.html'. This essentially
	// gives us a slice of all the 'page' templates for the application.
	pages, err := fs.Glob(ui.Files, "html/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Go through the template files from each page.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Create a slice containing the filepath patters for the templates we want to parse.
		patterns := []string{
			"html/base.layout.tmpl",
			"html/*.tmpl",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the resulting set of templates to the cache
		// using the page name 'home.page.tmpl' as a key for our map.
		cache[name] = ts
	}

	// Return the map we received.
	return cache, nil
}
