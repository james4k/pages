package pages

import (
	"net/http"
)

var DefaultGroup = New("pages", "pages/layouts")

// Dynamic returns an http.Handler with the named page loaded.
func Dynamic(name string, h Handler) http.Handler {
	return DefaultGroup.Dynamic(name, h)
}

// Static returns an http.Handler which serves the named static page. Panics on
// render error (while caching).
func Static(name string) http.Handler {
	return DefaultGroup.Static(name)
}
