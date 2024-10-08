package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/pkg/models"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// Define struct application for storing the dependencies of app.
// Field templateCache in dependencies struct, allow access to cache in all handlers.
// Handlers are in the same package, so we can define the functions as method against this struct, for access to the loggers.
// Add a sessionManager field to the application struct.
// Add a new users field to the application struct.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippeModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// New flag CLI, by default ":4000", info about flag, value of flag will save in addr variable.
	addr := flag.String("addr", ":4000", "network adress HTTP")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Function flag.Parse() for extract flag from CLI,
	// function reads flag's value from CLI and assign variable's content.
	// Have to call Parse function before use addr variable, otherwise it always will contain value by default ":4000",
	// if we have a mistakes while data is extracting, then our app break.
	flag.Parse()

	// Logger for record msgs about info errors with using stderr as a place for record.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Flag log.Lshortfile for log
	// file name and string number where errors were found
	// log.New() is safe for concurency using, we can share one logger for several Goroutines,
	// if we have several loggers and we use only one place for writing we have to be sure that method Write() also is safe for concurency using.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Connections in func openDB(), feed in func datasource(dsn) from flag cli.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use for closing pool of connections before the func main() is closed.
	defer db.Close()

	// Template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Decoder instance...
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager.
	// Then configure it to use or MySQL database as the session store,
	// set a lifetime of 12 hours, that sessions automatically expire
	// 12 hours after first being created.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.Secure = true

	// Structure with dependency injection.
	// Initialize a models.UserModel instance and add it to the application dependencies.
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippeModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this cae the only thing that we are changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Set the server's TLSConfig field to use the tlsConfig variable just created.
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Information in terminal about server launching.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// This func openDB covers sql.Open() and returns pool of connections sql.DB for current string of DSN connection.
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
