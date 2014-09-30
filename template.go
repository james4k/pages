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

func (t *Template) load(g *Group, name string, info interface{}) error {
	t.g = g
	tmpl := template.New(name)
	tmpl.Funcs(g.Funcs)
	var fm map[string]interface{}
	bytes, err := fmatter.ReadFile(filepath.Join(g.Dir, name), &fm)
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
	//writeInfo(info, fm)
	return nil
}

// writeInfo sets exported struct fields based on parsed frontmatter, in
// the form of a map. dest should be a pointer to a struct
func writeInfo(dest interface{}, fm map[string]interface{}) {
}

func (t *Template) Render(w io.Writer, data interface{}) error {
	t.g.mu.Lock()
	defer t.g.mu.Unlock()
	if !t.g.precache {
		t.g.layouts.Clear()
		err := t.g.layouts.Glob("*.html")
		if err != nil {
			panic(err)
			return err
		}
		err = t.load(t.g, t.tmpl.Name(), nil)
		if err != nil {
			panic(err)
			return err
		}
	}
	t.g.layouts.Funcs(t.funcMap)
	err := t.g.layouts.Execute(w, t.layout, t.tmpl, data)
	if err != nil {
		panic(err)
	}
	return nil
}
