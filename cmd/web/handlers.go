package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox/pkg/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the
	// manual check of r.URL.Path != "/" from this handler

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context.
	// We can use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// call method Get from for getting data by snippet's ID, if cant find snippet, then returns answer 404 error
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// use helper render() for display template
	app.render(w, r, "show.page.tmpl", data)
}

/* createSnippet page handler
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) //use method Header().Set() for add header 'Allow: POST' in map of http-headers, first parameter name of header, second value of header
		//w.WriteHeader(405)                       // we can call in handler only one time, for second time GO will give error for us. We have to call writeheader once before write for another status(instead 200 OK)
		//w.Write([]byte("Get method forbidden!\n"))
		app.clientError(w, http.StatusMethodNotAllowed) // using clientError() in helpers.go // old code: "http.Error(w, "Method is forbidden!", http.StatusMethodNotAllowed)" //we use func http.Error() for send different statuses
		return
	}
}
*/

// Add a new snippetCreate handler, which for now returns a placeholder response.
// We will update this shortly to show a HTML form.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Display the form for creating a new snippet..."))
	data := app.newTemplateData(r)

	app.render(w, r, "create.page.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the sam way for PUT and PATCH requests.
	// If there are any errors, we use our app. ClienError() helper to
	// send a 400 Bad Request response to the user
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*
	// However, we are expecting our expires value to be a number and what to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad Request
	// response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
