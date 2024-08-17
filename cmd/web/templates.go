package main

import "snippetbox/pkg/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet // add field Snippets in struct templateData
}

// create new type templateData which works as a storage for dynamic data that need to be passed to HTML templates
