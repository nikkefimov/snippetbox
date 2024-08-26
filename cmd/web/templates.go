package main

import (
	"html/template"
	"path/filepath"
	"snippetbox/pkg/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// initialize a new map which keeps cache
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// use the filepath.Glob() function to get a slice of all filepaths with '.page.tmpl'
	// essentially we will get a list of all the template files for the pages our application
	pages, err := filepath.Glob(filepath.Join(dir, "*page.tmpl"))
	if err != nil {
		return nil, err
	}

	// go through the template files from each page
	for _, page := range pages {
		// extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// process the iterate template file
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// use method ParseGlob for adding all framework templates
		// in our case it is only 'base.layout.tmpl' file
		ts, err = ts.ParseGlob(filepath.Join(dir, "*layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// use method ParseGlob to add all others templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// add the resulting set of templates to the cache
		// using the page name 'home.page.tmpl' as a key for our map
		cache[name] = ts
	}

	// return the map we received
	return cache, nil
}
