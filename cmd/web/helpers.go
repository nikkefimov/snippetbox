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

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Helper serverError writes error msg in errorLog and send to user answer 500 "Internal server error"
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Helper clientError sends exact status code and exact discription to user, in next steps, later it will look like 400 "Bad request"
// It happens when we have problem with user's request
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Helper notFound its a something like convenient shell around clientError, which sends to user answer "404 error"
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// extract the appropriate set of templates from the cache depending on the page name
	// if there is no entry of the requested template in the cache, then call the serverError() helper method
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s doesn't exist", name))
		return
	}

	// Initialize a new buffer
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
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst is the target
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
		// If we try to use a invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.
		// Use error.As() to check for this and raise a panic rather than returning the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil

}
