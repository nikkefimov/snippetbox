package main

import (
	"errors"
	"net/http"
	"strconv"

	"snippetbox/pkg/models"
)

// main page handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use render() for display template
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// showSnippet page handler
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // old code: "http.NotFound(w, r)", using notFound() in helpers.go
		return
	}

	// call method Get from for getting data by snippet's ID, if cant find snippet, then returns answer 404 error
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// use helper render() for display template
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// createSnippet page handler
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) //use method Header().Set() for add header 'Allow: POST' in map of http-headers, first parameter name of header, second value of header
		//w.WriteHeader(405)                       // we can call in handler only one time, for second time GO will give error for us. We have to call writeheader once before write for another status(instead 200 OK)
		//w.Write([]byte("Get method forbidden!\n"))
		app.clientError(w, http.StatusMethodNotAllowed) // using clientError() in helpers.go // old code: "http.Error(w, "Method is forbidden!", http.StatusMethodNotAllowed)" //we use func http.Error() for send different statuses
		return
	}
}
