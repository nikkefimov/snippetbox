package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox/pkg/models"
	"snippetbox/pkg/validator"

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

	// Call method Get from for getting data by snippet's ID, if cant find snippet, then returns answer 404 error
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

	// Use helper render() for display template
	app.render(w, r, "show.page.tmpl", data)
}

// Add a new snippetCreate handler, which for now returns a placeholder response.
// Will update this shortly to show a HTML form.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a greate opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, "create.page.tmpl", data)
}

// Remove the explicit FieldErrors struct field and instead embed the Validator type.
// Embediing this means that our snippetCreateForm 'inherits' all the
// fields and methods of our Validator type, including the FieldErrors field.
type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
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

	// Get the expires value from the form as normal.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create an instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		// Remove the FieldErrors assignment from here.
	}

	// Because the Validator type is embedded by the snippetCreateForm struct,
	// we can call CheckField() directly on it to execute our validation checks.
	// CheckField() will add the provided key and error message to he
	// FieldErros map if the check does not evaluate to true. For example, in
	// the first line here we "check that the form. Title field is not blank".
	// In the second, we "check that the form. Title field has a maximum character
	// length of 100" and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "title", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the Valid() method to see if any of the checks failed. If they did
	// then re-render the template passing in the form in the same way as before
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, "create.page.tmpl", data)
		return
	}

	// Also need to update this line to pass the data from the
	// snippetCreateForm instance to our Insert() method
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
