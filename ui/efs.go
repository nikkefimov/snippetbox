package ui

import "embed"

//go:embed "html" "static"
var Files embed.FS

// Important line here "//go:embed "html" "static"" looks like a comment, but it is actually a special comment directive.
// When application is compiled, this comment directive instructs Go to store the files from ui/html and ui/static folders in an
// embed.FS embedded filesystem referenced by the global variable Files.
