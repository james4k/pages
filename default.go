package pages

import "net/http"

var DefaultGroup = &Group{
	Dir:        "pages",
	LayoutsDir: "layouts",
}

/*
// Funcs adds template funcs to all pages and layouts that are loaded. See
// template.Funcs in html/template. This must be called before pages are loaded
// via Static or Dynamic.
func Funcs(f template.FuncMap) {
	DefaultGroup.Funcs(f)
}
*/

// SetPrecache enables precaching of templates for production usage. By
// default, precache is disabled to ease website development.
func SetPrecache(precache bool) {
	DefaultGroup.SetPrecache(precache)
}

func Parse(name string) (*Template, error) {
	return DefaultGroup.Parse(name)
}

func MustParse(name string) *Template {
	return DefaultGroup.MustParse(name)
}

// Handler returns an http.Handler which serves the named static page.
func Handler(name string, data interface{}) http.Handler {
	return DefaultGroup.Handler(name, data)
}
