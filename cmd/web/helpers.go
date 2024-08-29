package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// helper serverError writes error msg in errorLog and send to user answer 500 "Internal server error"

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// helper clientError sends exact status code and exact discription to user, in next steps, later it will look like 400 "Bad request"
// it happens when we have problem with user's request

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// helper notFound its a something like convenient shell around clientError, which sends to user answer "404 error"

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// extract the appropriate set of templates from the cache depending on the page name
	// if there is no entry of the requested template in the cache, then call the serverError() helper method
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s doesn't exist", name))
		return
	}

	// initialize a new buffer
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
