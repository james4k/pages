package pages

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"j4k.co/layouts"
)

// Handler that a Dynamic page implements by embedding pages.Template. NOT
// *pages.Template, unless you want to provide a non-nil pointer.
// TODO: user friendly panic on nil ptr? any way we can enforce this with type system?
type Handler interface {
	http.Handler
	load(*Group, string, interface{}) error
}

// TODO: Is there a reason to add mutexes for safe concurrency? Everything
// should be called at init, in the same function, but...hmm.
type Group struct {
	inited   bool
	precache bool
	layouts  *layouts.Group
	dir      string
	funcs    template.FuncMap
	mu       sync.Mutex
}

// New returns a new Group given paths to the layouts and pages. All .html
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

func (g *Group) SetPrecache(precache bool) {
	g.precache = precache
}

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
	err := h.load(g, name, nil)
	if err != nil {
		panic(err)
	}
	// TODO: reload in dev mode, just return h if not
	return h
}

type staticHandler struct {
	Template
	data interface{}
}

func (s *staticHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := s.Render(w, nil)
	if err != nil {
		log.Println(err)
	}
}

// Static returns an http.Handler which serves the named static page. Panics on
// render error (while caching).
func (g *Group) Static(name string, data interface{}) http.Handler {
	g.lazyInit()
	// TODO: should render template once to check for errors; panic if so
	// dev/live mode is maybe another story
	sh := &staticHandler{
		data: data,
	}
	h := g.Dynamic(name, sh)
	err := sh.Render(ioutil.Discard, nil)
	if err != nil {
		panic(err)
	}
	return h
}

/*
type dirHandler struct {
	Template
	dir      string
	notfound http.Handler
}

func (d *dirHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := d.Render(w, data)
	if err != nil {
		log.Println(err)
	}
}

func (d *dirHandler) read() {

	h := g.Dynamic(name, dh)
	err := dh.Render(ioutil.Discard, data)
	if err != nil {
		panic(err)
	}
}

// Dir returns an http.Handler which serves pages from a directory,
// non-recursively. Panics on render error (while caching).
func (g *Group) Dir(dir string, notfound http.Handler) http.Handler {
	g.lazyInit()
	// TODO: should render template once to check for errors; panic if so
	// dev/live mode is maybe another story
	dh := &dirHandler{
		dir:      dir,
		notfound: notfound,
	}
	dh.read()
	return dh
}
*/
