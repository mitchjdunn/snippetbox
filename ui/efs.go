package ui

import (
	"embed"
)

// this is a directive comment, on compile it will take the
// static folder and embed it in the embedded filesystem of
// our binary.  Root of fs is the package in which the directive
// is in, so ui
//
//go:embed "static" "html"
var Files embed.FS
