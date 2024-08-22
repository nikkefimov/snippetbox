package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"snippetbox/pkg/models/mysql"
)

// Define struct application for storing the dependencies of app
// Field templateCache in dependencies struct, allow access to cache in all handlers
// Handlers are in the same package, so we can define the functions as method against this struct, for access to the loggers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippeModel
	templateCache map[string]*template.Template //
}

func main() {

	// creating new flag CLI, by default ":4000"
	// adding info about flag
	// value of flag will save in addr variable
	addr := flag.String("addr", ":4000", "network adress HTTP")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Name of MySQL datasource")

	// we call function flag.Parse() for extract flag from CLI
	// function reads flag's value from CLI and assign variable's content
	// we have to call Parse function before use addr variable, otherwise it always will contain value by default ":4000"
	// if we have a mistakes while data is extracting, then our app break
	flag.Parse()

	// creat logger for record msgs about info errors with using stderr as a place for record
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// using flag log.Lshortfile for log
	// file name and string number where errors were found
	// log.New() is safe for concurency using, we can share one logger for several Goroutines
	// if we have several loggers and we use only one place for writing we have to be sure that method Write() also is safe for concurency using
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// creat connections code in func openDB(), we feed in func datasource(dsn) from flag cli
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// use for closing pool of connections before the func main() is closed
	defer db.Close()

	// initialise new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// define a new structure with dependency injection
	// initialise mysql.SnippetModel and add it in dependencies
	// add template cache in dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippeModel{DB: db},
		templateCache: templateCache,
	}

	// start a web server listening on porn :4000, using the mux as a router from file routes.go

	// call method app.routes() from the file routes.go
	// also in struct were created 'Addr and 'Handler' for the same network address and routes
	// and field ErrorLog for using logger by our server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// information in terminal about server launching
	infoLog.Printf("Launching server on %s", *addr)

	// method for logger and errors
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// this func openDB covers sql.Open() and returns pool of connections sql.DB for current string of DSN connection
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
