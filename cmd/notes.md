22.07 finish with new loggers and start dependency injection

handlers.go function home() is still writing errors in standart logger GO
working on it, have to change this to errorLog logger
update home() and other handler functions in handler.go
made mistake in file type "./ui/html/home.page.tmpl:" "changed to .bak" for testing new logger, tested - OK

23.07 creat helpers.go move some code with error handling, and update handlers.go with new features helpers

in handler serverError() was used func debug.Stack() for get trace from stack fo current goroutine and add it in logger,
its good in a future work because there is full route and easy for fixing. In helper clientError() was used func http.StatusText() for automatic text generation about status HTTP, like a "Bad request". Was used special constants from net/http for code about status HTTP instead number msgs. In helper serverError() was used constant http.StatusInternalServerError instead 500, in helper notFound() was used constant http.StatusNotFound instead 404.
Information about constants: "pkg.go.dev/net/http#pkg-constants"

23.07 fix msg from helpers.go, fix serverError() use Output(), depth is 2 by default

now have information about exact string in code in whole project with a problem, before had just information about string in helpers.go which says about problem is

23.07 correction of a specially made error earlier in type of file for the logger test and errors

mv ui/html/home.page.bak ui/html/home.page.tmpl. Tested - OK

25.07 create new file routes.go and new method, move this part with routes from main.go

after small refactoring updated file main.go is doing: parsing runtime configuration settings for an application, making dependencies for handlers, starting http server.

30.07 install homebrew in terminal, install Java JDK, install MySQL in terminal, launch MySQL trough brew

create new database "snippetbox", create new table "snippets", create tests notes, create new user for web with limited rights. Tested - OK.

download and install SQL driver for Go language from github.com

file go.mod updated according with installed SQL driver

file go.sum was created after install SQL driver, this file contains cryptographic checksums representing the contents of the required packages. Unlike the go.mod file, the go.sum file is not intended to be edited, and you should not normally opet it, much less edit it. This file accomplishes two useful taks: If you run the go mod veify command from a terminal, Go will check if the checksums of the of the dowloaded packages on your computer match the entries in go.sum, so you can be sure that they have not beed changed.
If someone else need to dowload all the dependencies for the project by running the go mod dowload command, will get an error message if there is any mismatch between the dependencies being downloaded and the checksums in the file.

2.08 Creat connections in MySQL, add sql.Open() func

Data source name for second parametr in sql.Open() func we can find github.com/go-sql-driver/mysql#dsn-data-source-name.
File main.go was updated

3.08 creat MySQL model for work with a database in project

create new folder mysql and two new files .go in folder models and mysql

in file models.go we define types of top level data, which our database model will use and return.

file snippets.go contains code for work with notes with MySQL database, assign new type here SnippetModel

7.08 update file snippets.go

edit method SnippetModel.Insert(), create new snippet in table snippets and return new snippet's id
make SQL request and update code in snippets.go, use interface sql.Result which we get after execution DB.Exec(). 
We have two methods from sql.Result, LastInsertId() and RowsAffected(), not all driver support these methods, PostgeSQL doesnt work with LastInsertId(), have to check driver's manual before use

8.08 output snippet from the database by snippet's ID from the URL

edit method GET in file snippets.go

13.08 test display latest snippets from DB

17.08 display content from MySQL into HTML template

fix error with mainpage

some work with templates:
-transfer dynamic data to HTML templates with scalable and secure way
-use various operators and functions from html/template package to control the display of dynamic data in a website template -cache the template so that resources are not wasted on re-processing the template for each HTTP request
-handle template rendering errors that occur at runtime
-realise a way to pass global dynamic data to web pages without reapeating code in handlers
-create custom functions to format and display data in HTML templates

-protection against XSS attacks in Go - data shielding
package "html/template" automatically escapes(screens) any data between {{}} tags, this behaviour helps to avoid cross-site scripting (XSS) attacks and is why use the "html/template" package instead of the simpler "text/template" package, also this package always removes any HTML comments you leave in template, including any conditional comments that are ofthen frontend developers make, this helps to avoid XSS attacks when dynamic content is displaying.

was created new template file "show.page.tmpl" and new file "templates.go" which contains new struct
tested - OK

work with operators and functions from Go template builder, was used {{define}}, {{template}}, {{block}}, {{if}}, {{with}}, {{range}}

updated template files for main page, tested - OK

17.08 template caching in Go
avoid processes the template files using the template.ParseFiles() function everytime when a webpage is displayed by processing the files once druing application startup and storing the processed templates in a cache in memory

put a code which reapets in handlers home and showSnippet in helper function

for caching processed templates using map in templates.go

initialise cache in main func

create new method render in helpers.go

update code in handlers.go for home() and showSnippet()

25.08 check whole code and comments

check errors with display single snippets

26.08 deliberate error 

test errors, add errors catcher in file helpers.go

28.08 

common dynamic data, updated footer

custom template functions

29.08 

-Middleware

create middleware.go and update routes.go

for check middleware info use curl request with a flag "curl -I http://localhost:4000/"

-Request logging

create logRequest() method using the standart middleware pattern

update middleware.go and routes.go

-Panic recovery, in a simple Go application when your code panics it will result in the application being terminated straight away. But in our application is a bit more sophisticated, Go HTTP server assumes that he effect of any panis is isolated to the goroutine serving the active HTTP request(every request is handled in its own goroutine).

if create deliberate panic in handlers.go, check by curl request, it would be Empty replry from server and empty response due to Go closing the underlying HTTP connection following the panic
this is not a greate experience for the user, it would be more appropriate and meaningful to send them a prope HTTP repsonse with a 500 Internal Server Error status instead

a neat way of doing this is to create some middleware which recovers the panic and calls our app.serverError() helper method

-Composable middleware chains, use justinas/alice package to help us manage our middleware/handler chains
its easy to create composable, reusable, middleware chains and that can be a real help application to grows and routes become more complex, the packgae itself is also small and lightweight and the code is clean and well written

update file routes.go with new package "Alice"

-Advanced routing, work with createSnippet handler
For GET /snippet/create requests adding a new snippet with a HTML form
For POST /snippet/create requests process this form data and then insert a new snippet record into database

Method Pattern              Handler             Action
GET    /                    home                Display the home page
GET    /snippet/view/:id    showSnippet         Display a specific snippet
GET    /snippet/create      createSnippet       Display a HTML for for creating a new snippet
POST   /snippet/create      createSnippetPost   Create a new snippet
GET    /static/             http.FileServer     Serve a specific static file

For some reasons Go's servemux doesnt support method based routing or clean URLs with variables in them,
most people tend to decide that its easier to reach for a third-party package to help with routing (julienschmidt/httprouter, go-chi/chi and gorilla/mux) this all three support method-based routing and clean URLs, but beyond that they have lightly different behaviours and features.

In summary:
*julienschmidt/httprouter is the most focused, lightweight and fastest of the three packages, and is about as close to 'perfect' as any third-party router gets in terms of its compliance with the HTTP specs. It automaticly handles OPTIONS requests and sends 405 responses correctly, and allows you to set custom handlers for 404 and 405 responses too.
*go-chi/chi is generally similar to httprouter in terms of its featues, with the main differences being that it also supports regexp route patterns and 'grouping' of routes which use specific middleware. This route grouping features is really valuable in larger applications where you have lots routes and middleware to manage(two downsides of chi are that it doesnt automatically handle OPTIONS requests and it doesnt set an Allow header in 405 responses).
*gorilla/mux is the most full-featured of the tree routers. It supports regexp route patterns, and allows to route requests based on scheme, host and headers. Its also the only one to support custom routing rules and route 'reversing'(like you get in Django, Rails, or Laravel). The main downside of gorilla/mux is that its comparatively slow and memory hungry - although for a dayabase-driven web application like ours app the impact over the lifetime of a whole HTTP request is likely to be small. Like chi, it also doesnt automatically handle OPTIONS requests and it doesnt set an Allow header in 405 responses.

In our case, our application is fairly small and we dont need support for anuthing beyond basic method-based routing and clean URLs. So, for the sake of performance and correctness, we will opt to use julienschmidt/httprouter in this project.

-Clean URLs and method-based routing

install httprouter package, update routes.go and handlers.go and template home.page

30.08 Processing forms

-Update handlers.go, create createSnippet handler, update create.page.tmpl

31.08 Parsing for data

-Update snippetCreatePost method in handlers.go file.

-Accessed the form values via the r.PostForm map. But an alternative approach is to use the(subtly different) r.Form map.
The r.PostForm pas is populated only for POST, PATCH and PUT requests, and contains the form data from the request body.

In contrast, the r.Form map is populated for all requests(irrespective of their HTTP method), and contrain the form data from any request body and any query string parameters. So if our form was submitted to /snippet/create?foo=bar, we could also het the value of the foo parameter by calling r.Form.Get("foo"). Note that in the event of a conflict, the request body value will take precedent over the query string parameter.

Using the f.Form map can be useful if your application sends data in a HTML form and in the URL, or you have an application that is agnostic about how parameters are passed. Our case this things are not applicable, expect our form data to be sent in the request body only, so it is for sensible for us to acces it via r.PostForm

1.09 Validation form data
 
 -Update handlers.go, create a map

 When we check the length of the title field we are using the utf8.RuneCountInString()function - notGo's len() function. This is because we want to count the number of characters in the title rather than the number of bytes. To illustrate the difference, the sting "šop" has 3 characters but a length of 4 bytes because of the umlauted š character.

 -Displaying erros and repopulating fields

 update snippetCreatePost and create.page template file (unlike struct fields, map key names dont have to be capitalized in order to access them from a template)

 For the validation errors, the underlying type of our FieldErrors field is a map[string]string, which uses the form field names as keys. For maps, its possible to access the value for a given key by simply chaining the key name. So, for example, to render a validation error for the title field we can use the tag{{.Form.FieldErrors.title}} in our template.

 -Creating validation helpers

 Update code in handlers.go and create validator.go (validator package)

 -Automatic form parsing

 download package goplayground/form or gorila/schema to automatically decode the form data into the crateSnippetForm struct

add package and update files main.go and handlers.go

When call app.formDecoder.Decode() it requires a non-nil pointer as the target decode destination. If we try to pass in something that is not a non-nil pointer, then Decode() will return a form.InvalidDecodeError error.
It is a critical problem with our application code(rather than a client error due to bad input). Need to check for this error specifically and manage it as a special case, rather than just returning a 400 Bad Request response.

Creating a decodePostForm helper, update helpers.go file

02.09 Stateful HTTP

A confirmation message like this should only show up for the user once (after creating snippet) and no other user should ever see the message.

There are a lot of security considerations when it comes to working with sessions and proper implementation is not trivial. Unless you really need to roll your own implementation, its a good idea to use an existing, well-tested, third-party package.

*gorilla/sessions
is the most establishe dand well-known sessian management package for Go. It has a simple and easy to use API and let's you stroe session data clien side(in signed and ecrypted cookies) or server side(in a database like MySQL), PostgeSQL or Redis.
It doesnt provide mechanism to renew session IDs, which is necessary to reduce risks associated with session fixation attacks if you are using one of the server side session stores.

*alexedwards/scs 
lets store session data server side only. It supports automatic loading and saving of session data via middleware, has a nice interface fpr type safe manipulation of data and does allow renvewal of session IDs. Like gorilla/sessions, it also supports a variety of databases including MySQL, PostgreSQL and Redis.

In summary, if you want to store session data client side in a coockie then gorilla/sessions is a good choice, but otherwise alexedwards/scs is generally the better option due to the ability to renew sessionn IDs.

get github.com/alexedwards/scs/v2@v2
get github.com/alexedwards/scs/mysqlstore

-Setting up the session manager

use alexedwards/scs package, before need to do create a sessions table in MySQL database to hold the session data for users.

update main.go
The scs.New() function returns a pointer to a SessionManager struct which holds configuration settings for your sessions. In the code have set the Store and Lifetime fields of this struct, but there is a range of other fields that you can and should configure depending on application's need.

update routes.go

-Working with session data
set the session functionality to work and use it to persist the confirmation flash message between HTTP requests.
Update handlers.go, templates.go, base.layout.tmpl

-Auto displaying flash messages
automate the display of flash messages, that any message is automatically included the next time any page is rendered.
That change means that no longer need to check for the flash message within the snippetView handler.

update helpers.go, edit handlers.go

03.08 Security inprovements

Make some improvements to application, secure data during transit and make server better to some common types of denial of service attacks.

-Generating a self-signed TLS certificate.

For MacOS, FreeBSD or Linux the generate_cert.go file should be located under: "/usr/local/go/crypto/tls/"

For generate certificate execute in terminal in folder tls: "go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost"

First it generates a 2048-bit RSA key pair, which is a cryptographically secure public key and private key. It then stroes the private key in a key.pem file and generates a self signed TLS certificate for the host localhost containing the public key, which is stores in a cert.pem. Both the private key and certificate are PEM encoded, which is the standart format used by most TLS implementations.

Now application has a self signed TLS certificate and corresponding private key that can be use during development.

04.08

-Running a HTTPS server

Now starting a HTTPS web server, just need make some changes in main.go and swap srv.ListenAndServe() to swap srv.ListenAndServe() instead.
After that, the only difference is that it will now be talking HTTPS instead of HTTP (https://localhost:4000/)
Application homepage should appear (although it will still carry a warning in the URL bar because the TLS certificate is self-signed).

If inspect page, we will see in security technical details section, that connection is encrypted and working as expected.
That TLS version 1.3 is being used and che cipher suite for HTTPS connection is TLS_AES_128_GCM_SHA256.

Important to know that HTTPS server only supports HTTPS. If try making a regular HTTP request to it, the server will send the user a 400 Bad Request status and the message "CLient sent an HTTP request to an HTTPS server". Nothing will be logged.

A big plus of using HTTPS is that, if a client supports HTTP/2 connections - Go's HTTPS server automatically upgrade the connection to use HTTP/2. It's good because it means, that ultimately our pages will load faster for users. 

Important to note that the user that using to run Go application must have read permissions for both the cert.pem and key.pem files, otherwise ListenAndServeTLS() will return a permision denied error. By default the generate_cert.go tool grants read permission to all users for the cert.pem file, but read permission only to the owner of the key.pem file.

For version control system, may to add an ignore rule by "eco 'tls/' >> .gitignore"

-Configure HTTPS settings

Go has good fedault settings for HTTPS server, but it is possible to optimize and customize how the server behaves.
Go support a few elliptic curves. but as of Go 1.18 only tls.CurveP256 and tls.X25519 have assembly implementations. The others are very CPU intensive, so omitting them helps ensure that our server will remain performant under heave loads. To make this tweak, create a tls.Config struct containing our non default TLS settings and add it to our http.Server struct before we start the server.

update main.go

TLS versions are also defined as constants in the crypto/tls package and Go's HTTPS server supports TLS versions 1.0 to 1.3.
You can configure the minimum and maximum TLS versions via the tls.Config.MinVersion and MaxVersion fields. For instance, that all computers in your user base support TLS 1.2, but not TLS 1.3, then you may wish to use a cinfiguration like so:
"tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12.
    MacVersion: tls.VersionTLS12,
}"

-Restricting cipher suites
The full range of cipher suites that Go supports are defined in the crypto/tls package constants.

For some applications, it may be be desirable to limit your HTTPS server to only support some of these cipher suites. For example you might want to only support cipher suites which use ECDHE(forward secrecy) and not support weak cipher suites that use RC4, 3DES or CBC. You can do this via the tls.COnfig.CupherSuites field like so:

tlsConfig := &tls.Config{
    CipherSuites: []uint16{
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
    },
}

Go will automatically choose which of these cipher suites is actually used at runtime based on the cipher security, perfomance and client/server hardware support.

Important to know, restricting the supported cipher suites to only include strong, modern, ciphers can mean that users with certain older browsers will not able to use your website. There is a balance to be struck between security and backwards compatibility and the right decision for you will depend on the technology typically used by your user base.
Its also important to note that if a TLS 1.3 connection is negotiated any CipherSuites field in your tls.Config will be ignored. The reason for this is that all the cipher suites that Go supports fot TLS 1.3 connections are vonsidered to be a safe, so there is not much point in providing a mechanism to configure them. Basically, using tls.Config to set a custom list of supported cipher suites will affect TLS 1.0-1.2 connections only.

-Connection timeouts

Improve the resilency of server by adding some timeout settings: "IdleTimeout, ReadTimeout, WriteTimeout".
All three of these timeouts - IdleTimeout, ReadTimeout and WriteTimeout - are server-wide settings which act on the underlying connection and apply to all requests irrespective of their handler or URL.

*The IdleTimeout setting.
By default, Go enable keep-alives on all accepted connections. This helps reduce latency because a client can reuse the same connection for multiple requests without having to reapeat the handshake.
By default, keep-alive connections will be automatically closed after a couple of minutes (depending on your OS). This helps to clear-up connections where the user has unexpectedly disappeared - e.g. due to a power cut cliend-side.
There is no way to increase this default (unless you roll your own net.Listener), but you can reduce it via the IdleTimeout setting. In our case, we have set IdleTimeout to 1 minute, which means that all keep-alive connections will be automatically closed after 1 minute of inactivity.

*The ReadTimetout setting.
In code we have also set the ReadTimeout setting to 5 seconds. This means that if the request header or body are still being read 5 second after the request is first accepted, the Go will close the underlying connection. Because this is a 'hard' closure on the connection the user will not receive any HTTP(S) response. Setting a short ReadTimeout period helps to mitigate the risk from slow-client attacks suck as Slowloris - which could otherwise keep a connection open endefinitely by sending partial incomplete, HTTP(S) requests. (If you set ReadTimeout bud dont set IdleTimeout, then IdleTimeout will default to using the same setting as ReadTimeout). Generally, recommendation is to avoyd any ambiguity and always set an explicit IdleTimeout value for your server.

*The WriteTimeout setting.
The WriteTimeout setting will close the underlying connection if our server attempts to write to the connection after a given period(in our code its 10 seconds). But this behaves slightly defferently depending on the protocol being used.
For HTTP, if some data is written to the connection more than 10 seconds after read of the request header finished, Go will close the underlying connection instead of writing thedata.
For HTTPS connections, if some data is written to the connection mroe than 10 seconds after the request is first accpeted, Go will close the underlying connection instead of writing the data. This means that if you are using HTTPS its sensible to set WriteTimeout to a value greater than ReadTimeout.
Therefore, the idea of WriteTimeout is generally not to prevent long-running handlers, but to prevent the data that the handler returns from taking too long to write.

*The MaxHeaderBytes setting.
The http.Server object also provides a MaxHeaderBytes field, which you can use to control the maximum number of bytes the server will read when parsing request headers. Bydefualt, Go allows a maximum header length of 1 MB. To change (for 0.5MB) this:
" srv := &http.Server{
    Addr: *addr,
    MaxHeaderBytes: 524288,
}"

If MaxHeaderBytes is exceeded then the user will automatically be sent a 431 Request Header Fields Too Large response.

update main.go

-User authentication

*A user will register by visiting a form at /user/signup and entering name, email, address and password. This information in a new users database table.
*A user will log in by visiting a form at /user/login and entering email address and password.
*Then check the dayabase to see if the email and password entered match one of the users in the users table. If there is a match, user has authenticated successfully and we add the relevant id value for the user to their session data, using the key "authenticatedUserID".
*When we receive any subsequent requests, we can check the user's session daya for a "authenticatedUserID" value. If it exists, we know that the user has already successfully logged in. We can keep checking this until the session expires, when the user will need to log in again. If there is no "authenticatedUserID" in the session, we know that the user is not logged in.

-Routes setup
Method  Pattern            Handler           Action
GET     /                  home              Display    the home page
GET     /snippet/view/:id  showSnippet       Display a specific snippet
GET     /snippet/create    snippetCreate     Display a HTML form for creating a new snippet
POST    /snippet/create    snippetCreatePost Create a new snippet
GET     /user/signup       userSignup        Display a HTML form for signing up a new user
POST    /user/signup       userSignupPost    Create a new user
GET     /user/login        userLogin         Display a HTML form for logging in a user
POST    /user/login        userLoginPost     Authenticate and login the user
POST    /user/logout       userLogoutPost    Logout the user
GET     /static/*filepath  http.FileServer   Serve a specific static file

updata handlers.go and routes.go

-Creating a users model
Set up a new users database table and a database model to access it.

*Creating new sql table 'users' in terminal

"Use snippetbox;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email); "

The id field is an autoincrementing integer field and the primary key for the table, This means that the users ID values are guaranteed to be unique positive integers.
The type of hashed_password field, because storing hashes of the user password in the database, not the password themselves and the hashed versions will always be exactly 60 characters long.
Added a UNIQUE constraint on the email column, will not up with two users who have the same email address.

*Building the model in Go
update errors.go, create users.go, update main.go