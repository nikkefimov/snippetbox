package main

import (
	"html/template"
	"path/filepath"
	"snippetbox/pkg/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet // add field Snippets in struct templateData
}

// create new type templateData which works as a storage for dynamic data that need to be passed to HTML templates

func newTemplateCache(dir string) (map[string]*template.Template, error) { // new map for store cache
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// use func filepath.Glob for get slice with file routs with a type 'page.tmpl'

	for _, page := range pages { // going through the template files from each page
		name := filepath.Base(page)          // exctract part of file's name '.tmpl' and assigning it to the name variable
		ts, err := template.ParseFiles(page) // process the iterated template file
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl")) // use method ParseGlob for adding all framework templates, like a file 'base.layout.tmpl'
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*partial.tmpl")) // use method ParseGlob for adding additional templates, like a footer.partial.tmpl
		if err != nil {
			return nil, err
		}
		cache[name] = ts // add the resulting set of templates to the cache using the page name as a key for our map
	}

	return cache, nil // return the card we received
}
