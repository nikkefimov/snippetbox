package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w) // old code "http.NotFound(w, r)", using notFound() in helpers.go
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{ // use render() for display template
		Snippets: s,
	})

	/*for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}*/

	data := &templateData{Snippets: s} // create example of templateData struct which contains slice with snippets

	// checking unexist pages

	files := []string{ //creating slice which contains route for two tmpl files, file home.page.tmpl must go first in list
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...) // using func template.ParseFiles() for read our template
	if err != nil {                          // if error we write specify msg about error and use func http.Error() for send this info to user
		app.serverError(w, err) // updated, using serverError() in helpers.go
		return
	}

	// transfer struct templateData in template and now struct is available inside of template files by using . before name
	err = ts.Execute(w, data) //we use func Execute() for write template's content in body of http response. Last parameter in Execute func needs for send dynamic data in template
	if err != nil {
		app.serverError(w, err) // using serverError() in helpers.go
	}
}

// main page

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // old code: "http.NotFound(w, r)", using notFound() in helpers.go
		return
	}

	s, err := app.snippets.Get(id) // call method Get from model for get data by snippet's ID, if cant find snippet, than returns answer 404
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{ // use helper render() for display template
		Snippet: s,
	})
	// fmt.Fprintf(w, "%v", s) // show all returns on the page

	data := &templateData{Snippet: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}
	//initialise slice contains rout to file show.page.tmpl, add base template and part of footer with allready made earlier

	ts, err := template.ParseFiles(files...) //parsing files
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data) // execute snippet with data, transfer the templateData structure as the data for the template
	if err != nil {
		app.serverError(w, err)
	}
}

// display notes

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) //use method Header().Set() for add header 'Allow: POST' in map of http-headers, first parameter name of header, second value of header
		//w.WriteHeader(405)                       // we can call in handler only one time, for second time GO will give error for us. We have to call writeheader once before write for another status(instead 200 OK)
		//w.Write([]byte("Get method forbidden!\n"))
		app.clientError(w, http.StatusMethodNotAllowed) // using clientError() in helpers.go // old code: "http.Error(w, "Method is forbidden!", http.StatusMethodNotAllowed)" //we use func http.Error() for send different statuses
		return
	}

	title := "First story"
	content := "about first story, first"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires) // transfet data in method SnippetModel.Insert(), and taking back ID of the newly created record into the database
	if err != nil {
		app.serverError(w, err)
		return
	}
	//w.Write([]byte("form for creating note..."))
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther) // redirect user to page with note ID
}

//use r.Method for check type of request, error only for method GET
//notes handler
