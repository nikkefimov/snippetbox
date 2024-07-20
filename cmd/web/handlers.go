package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// checking unexist pages

	files := []string{ //creating slice which contains route for two tmpl files, file home.page.tmpl must go first in list
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...) //using func template.ParseFiles() for read our template
	if err != nil {                          //if error we write specify msg about error and use func http.Error() for send this info to user
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500) //msg about server error(inside)
		return
	}

	err = ts.Execute(w, nil) //we use func Execute() for write template's content in body of http response. Last parameter in Execute func needs for send dynamic data in template
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	//w.Write([]byte("Hello from Snippetbox"))
}

// main page

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// it was before changes w.Write([]byte("showing note..."))
	fmt.Fprintf(w, "Show selected note with ID %d...", id) //as a first parametr we used w, instead io.Writer
}

// display notes

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) //use method Header().Set() for add header 'Allow: POST' in map of http-headers, first parameter name of header, second value of header
		//w.WriteHeader(405)                       // we can call in handler only one time, for second time GO will give error for us. We have to call writeheader once before write for another status(instead 200 OK)
		//w.Write([]byte("Get method forbidden!\n"))
		http.Error(w, "Method is forbidden!", http.StatusMethodNotAllowed) //we use func http.Error() for send different statuses
	}
	w.Write([]byte("form for creating note..."))
}

//use r.Method for check type of request, error only for method GET
//notes handler
