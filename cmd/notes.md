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

2.07 Creat connections in MySQL, add sql.Open() func

Data source name for second parametr in sql.Open() func we can find github.com/go-sql-driver/mysql#dsn-data-source-name.
File main.go was updated

3.07 creat MySQL model for work with a database in project

create new folder mysql and two new files .go in folder models and mysql

in file models.go we define types of top level data, which our database model will use and return.

file snippets.go contains code for work with notes with MySQL database, assign new type here SnippetModel

7.07 update file snippets.go

edit method SnippetModel.Insert(), create new snippet in table snippets and return new snippet's id
make SQL request and update code in snippets.go, use interface sql.Result which we get after execution DB.Exec(). 
We have two methods from sql.Result, LastInsertId() and RowsAffected(), not all driver support these methods, PostgeSQL doesnt work with LastInsertId(), have to check driver's manual before use

8.07 output snippet from the database by snippet's ID from the URL

edit method GET in file snippets.go

13.07 test display latest snippets from DB

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

17.07 template caching in Go
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

