package main

import (
	"html/template"
	"path/filepath"
	"snippetbox/pkg/models"
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
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template, FuncMap object and store it in a global variable
// this is essentially a string-keyd map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Initialize a map which keeps cache.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths with '.page.tmpl',
	// essentially will get a list of all the template files for the pages our application.
	pages, err := filepath.Glob(filepath.Join(dir, "*page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Go through the template files from each page.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Process the iterate template file.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use method ParseGlob for adding all framework templates
		// in our case it is only 'base.layout.tmpl' file.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use method ParseGlob to add all others templates.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*partial.tmpl"))
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
