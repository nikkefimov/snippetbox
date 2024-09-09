package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// Helper serverError writes error msg in errorLog and send to user answer 500 "Internal server error".
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Helper clientError sends exact status code and exact discription to user, in next steps, later it will look like 400 "Bad request".
// It happens when we have problem with user's request.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Helper notFound it's a something like convenient shell around clientError, which sends to user answer "404 error".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Extract the appropriate set of templates from the cache depending on the page name,
// if there is no entry of the requested template in the cache, then call the serverError() helper method.
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s doesn't exist", name))
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// Add a flash message to the template data, if one exists.
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		// Add the authentication status to the template data.
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst is the target,
// destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same wat that we did in our createSnippetPost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		// If we try to use a invalid target destination, the Decode() method,
		// will return an error with the type *form.InvalidDecoderError.
		var invalidDecoderError *form.InvalidDecoderError

		// Use error.As() to check for this and raise a panic rather than returning the error.
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil

}

// Return true if the current request is from an authenticated user, otherwise return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
