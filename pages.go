package pages

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"j4k.co/layouts"
)

// TODO: Is there a reason to add mutexes for safe concurrency? Everything
// should be called at init, in the same function, but...hmm.
type Group struct {
	inited   bool
	precache bool
	layouts  *layouts.Group
	mu       sync.Mutex

	Dir        string
	LayoutsDir string
	Funcs      template.FuncMap
}

func (g *Group) lazyInit() {
	if g.inited {
		return
	}
	g.inited = true
	g.layouts = layouts.New(g.LayoutsDir)
	g.layouts.Funcs(g.Funcs)
	g.layouts.Funcs(template.FuncMap{
		"page": emptyPageFn,
	})
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

/*
// Funcs adds template funcs to all pages and layouts that are loaded. See
// template.Funcs in html/template. This must be called before pages are loaded
// via Static or Dynamic.
func (g *Group) funcs(f template.FuncMap) {
	if g.Funcs == nil {
		g.Funcs = template.FuncMap{}
	}
	for k, v := range f {
		g.Funcs[k] = v
	}
	g.layouts.Funcs(f)
	// To support late Funcs calls, we would need a registry of our page
	// templates. Not sure it's worth it.
}
*/

// New returns a Template loaded using the configured page/layout paths.
func (g *Group) Parse(name string) (*Template, error) {
	g.lazyInit()
	var t Template
	err := t.load(g, name, nil)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (g *Group) MustParse(name string) *Template {
	t, err := g.Parse(name)
	if err != nil {
		panic(err)
	}
	return t
}

type staticHandler struct {
	tmpl *Template
	data interface{}
}

func (s *staticHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := s.tmpl.Render(w, s.data)
	if err != nil {
		log.Println(err)
	}
}

// Handler returns an http.Handler which serves the named page.
func (g *Group) Handler(name string, data interface{}) http.Handler {
	g.lazyInit()
	// TODO: should render template once to check for errors; panic if so
	// dev/live mode is maybe another story
	h := &staticHandler{
		tmpl: g.MustParse(name),
		data: data,
	}
	err := h.tmpl.Render(ioutil.Discard, data)
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
