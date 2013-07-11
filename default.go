package pages

import (
	"html/template"
	"net/http"
)

var DefaultGroup = New("pages", "pages/layouts")

// NewDefault initializes a new default group given paths to the layouts and
// pages. All .html files in the layouts path are loaded
func NewDefault(pagesPath, layoutsPath string) {
	DefaultGroup = New(pagesPath, layoutsPath)
}

// Funcs adds template funcs to all pages and layouts that are loaded. See
// template.Funcs in html/template. This must be called before pages are loaded
// via Static or Dynamic.
func Funcs(f template.FuncMap) {
	DefaultGroup.Funcs(f)
}

// Dynamic returns an http.Handler with the named page loaded.
func Dynamic(name string, h Handler) http.Handler {
	return DefaultGroup.Dynamic(name, h)
}

// Static returns an http.Handler which serves the named static page. Panics on
// render error (while caching).
func Static(name string) http.Handler {
	return DefaultGroup.Static(name)
}
