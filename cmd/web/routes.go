package main

import (
	"net/http"
	"snippetbox/pkg/nfs"

	"github.com/justinas/alice"
)

// update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServerMux
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// use package nfs from nfs.go
	// use FileServer for processing http requests for static files from folder
	fs := http.FileServer(nfs.NeuteredFileSystem{Fs: http.Dir(",/ui/static")})
	mux.Handle("/static", http.StripPrefix("/static", fs))

	// use new package alice
	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standart.Then(mux)
}
