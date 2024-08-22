package main

import (
	"net/http"
	"snippetbox/pkg/nfs"
)

// use mux as a router, create method routes()
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// use package nfs from nfs.go
	// use FileServer for processing http requests for static files from folder
	fs := http.FileServer(nfs.NeuteredFileSystem{Fs: http.Dir(",/ui/static")})
	mux.Handle("/static", http.StripPrefix("/static", fs))

	return mux
}
