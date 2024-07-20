package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	//creating new flag CLI, by default ":4000"
	//adding info about flag
	//value of flag will save in addr variable

	addr := flag.String("addr", ":4000", "network adress HTTP")
	flag.Parse()

	//we call function flag.Parse() for extract flag from CLI
	//function reads flag's value from CLI and assign variable's content
	//we have to call Parse function before use addr variable, otherwise it always will contain value by default ":4000"
	//if we have a mistakes while data is extracting, then our app break

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//creating logger for record msgs about info errors with using stderr how place for record
	//using flag log.Lshortfile for log
	//file name and string number where errors was found
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")}) //use FileServer for processing http requests for static files from folder ./ui/static. http.Dir its root project folder
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{ // struct with server information and for new logger
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Launching server on %s", *addr)
	// old logger "err := http.ListenAndServe(*addr, mux)"
	err := srv.ListenAndServe() // new logger with new struct
	errorLog.Fatal(err)
	//we can redirect msg from terminal to log txt file on HDD with "go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log"
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
