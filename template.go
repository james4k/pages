package pages

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/russross/blackfriday"

	"j4k.co/fmatter"
)

type Template struct {
	g       *Group
	tmpl    *template.Template
	layout  string
	funcMap template.FuncMap
}

func emptyPageFn() interface{} {
	return map[string]interface{}{}
}

func (t *Template) Render(w io.Writer, data interface{}) error {
	// TODO: To improve throughput, we can probably use html/template.Clone()
	// to sort of bind a layouts template to a page with unique FuncMaps. So in
	// theory there will need to be no synchronization, and yet they will share
	// the same parse trees. I'm sure this will also require changes to
	// j4k.co/layouts. We might be able to also 'flatten' the recursion done by
	// layouts to a self executing html/template.
	t.g.mu.Lock()
	defer t.g.mu.Unlock()
	if !t.g.precache {
		t.g.layouts.Clear()
		err := t.g.layouts.Glob("*.html")
		if err != nil {
			return err
		}
		err = t.load(t.g, t.tmpl.Name(), nil)
		if err != nil {
			return err
		}
	}
	t.g.layouts.Funcs(t.funcMap)
	return t.g.layouts.Execute(w, t.layout, t.tmpl, data)
}

func (t *Template) load(g *Group, name string, info interface{}) error {
	t.g = g // maybe create a setGroup method instead
	tmpl := template.New(name)
	tmpl.Funcs(g.funcs)
	var fm map[string]interface{}
	bytes, err := fmatter.ReadFile(filepath.Join(g.dir, name), &fm)
	if err != nil {
		return err
	}
	t.funcMap = template.FuncMap{
		"page": func() interface{} {
			return fm
		},
	}
	tmpl.Funcs(t.funcMap)
	ext := filepath.Ext(name)
	if ext == ".md" {
		bytes = blackfriday.MarkdownBasic(bytes)
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
	writeInfo(info, fm)
	return nil
}

// writeInfo sets exported struct fields based on parsed frontmatter, in
// the form of a map. dest should be a pointer to a struct
func writeInfo(dest interface{}, fm map[string]interface{}) {
}
