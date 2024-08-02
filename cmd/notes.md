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