package pages

import (
	"html/template"
	"io"
	"io/ioutil"
	"j4k.co/fmatter"
	"j4k.co/layouts"
	"net/http"
	"path/filepath"
)

// Handler that a Dynamic page implements by embedding pages.Template. NOT
// *pages.Template, unless you want to provide a non-nil pointer.
// TODO: user friendly panic on nil ptr? any way we can enforce this with type system?
type Handler interface {
	http.Handler
	load(*Group, string) error
}

type Template struct {
	g      *Group
	tmpl   *template.Template
	layout string
}

func (t *Template) Render(w io.Writer, data interface{}) error {
	return t.g.layouts.Execute(w, t.layout, t.tmpl, data)
}

func (t *Template) load(g *Group, name string) error {
	t.g = g // maybe create a setGroup method instead
	tmpl := template.New(name)
	//tmpl.Funcs(g.funcs)
	var fm map[string]interface{}
	bytes, err := fmatter.ReadFile(filepath.Join(g.dir, name), &fm)
	if err != nil {
		return err
	}
	_, err = tmpl.Parse(string(bytes))
	if err != nil {
		return err
	}
	t.tmpl = tmpl
	if l, ok := fm["layout"]; ok {
		t.layout = l.(string)
	} else {
		t.layout = "default"
	}
	// TODO: find FrontMatter field in handler and unmarshal into that
	return nil
}

type Group struct {
	layouts *layouts.Group
	dir     string
}

// New returns a new Pages, given paths to the layouts and pages. All .html
// files in the layouts path are loaded.
func New(pagesPath, layoutsPath string) *Group {
	g := &Group{
		layouts: layouts.New(layoutsPath),
		dir:     pagesPath,
	}
	g.layouts.Glob("*.html")
	return g
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

// TODO: Funcs method

// Dynamic returns an http.Handler with the named page loaded into your
// embedded pages.Template.
func (g *Group) Dynamic(name string, h Handler) http.Handler {
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
	// TODO: should render template once to check for errors; panic if so
	// dev/live mode is maybe another story
	sh := &staticHandler{}
	h := g.Dynamic(name, sh)
	// check that it renders without error
	err := sh.Render(ioutil.Discard, sh.FrontMatter)
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
