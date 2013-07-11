package pages

import (
	"html/template"
	"io/ioutil"
	"j4k.co/layouts"
	"net/http"
	"sync"
)

// Handler that a Dynamic page implements by embedding pages.Template. NOT
// *pages.Template, unless you want to provide a non-nil pointer.
// TODO: user friendly panic on nil ptr? any way we can enforce this with type system?
type Handler interface {
	http.Handler
	load(*Group, string) error
}

// TODO: Is there a reason to add mutexes for safe concurrency?
type Group struct {
	inited  bool
	layouts *layouts.Group
	dir     string
	funcs   template.FuncMap
	mu      sync.Mutex
}

// New returns a new Pages, given paths to the layouts and pages. All .html
// files in the layouts path are loaded. Panics on error as common usage is
// assignment to package scoped variables.
func New(pagesPath, layoutsPath string) *Group {
	g := &Group{
		layouts: layouts.New(layoutsPath),
		dir:     pagesPath,
	}
	g.layouts.Funcs(template.FuncMap{
		"page": emptyPageFn,
	})
	return g
}

func (g *Group) lazyInit() {
	if g.inited {
		return
	}
	g.inited = true
	err := g.layouts.Glob("*.html")
	if err != nil {
		panic(err)
	}
}

/*
func (g *Group) SetPaths(pagesPath, layoutsPath string) {
	g.layouts.SetPath(layoutsPath)
	g.dir = pagesPath
}
*/

/*
func (g *Group) NoCache(nocache bool) {
}
*/

// Funcs adds template funcs to all pages and layouts that are loaded. See
// template.Funcs in html/template. This must be called before pages are loaded
// via Static or Dynamic.
func (g *Group) Funcs(f template.FuncMap) {
	if g.funcs == nil {
		g.funcs = template.FuncMap{}
	}
	for k, v := range f {
		g.funcs[k] = v
	}
	g.layouts.Funcs(f)
	// To support late Funcs calls, we would need a registry of our page
	// templates. Not sure it's worth it.
}

// Dynamic returns an http.Handler with the named page loaded into your
// embedded pages.Template.
func (g *Group) Dynamic(name string, h Handler) http.Handler {
	g.lazyInit()
	err := h.load(g, name)
	if err != nil {
		panic(err)
	}
	// TODO: reload in dev mode, just return h if not
	return h
}

// Static returns an http.Handler which serves the named static page. Panics on
// render error (while caching).
// TODO: pass in template data?
func (g *Group) Static(name string) http.Handler {
	g.lazyInit()
	// TODO: should render template once to check for errors; panic if so
	// dev/live mode is maybe another story
	sh := &staticHandler{}
	h := g.Dynamic(name, sh)
	err := sh.Render(ioutil.Discard, nil)
	if err != nil {
		panic(err)
	}
	return h
}

// TODO: cache result in memory, gzipped. Maybe..how much should we care about
// memory usage?
type staticHandler struct {
	Template
	FrontMatter map[string]interface{}
}

func (s *staticHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Render(w, s.FrontMatter)
}
