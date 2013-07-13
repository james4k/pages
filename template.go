package pages

import (
	"html/template"
	"io"
	"j4k.co/fmatter"
	"path/filepath"
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
		t.g.layouts.Glob("*.html")
		t.load(t.g, t.tmpl.Name())
	}
	t.g.layouts.Funcs(t.funcMap)
	return t.g.layouts.Execute(w, t.layout, t.tmpl, data)
}

func (t *Template) load(g *Group, name string) error {
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
	// TODO: find FrontMatter field in handler and unmarshal into that. At
	// least for the Dynamic handler.
	return nil
}
