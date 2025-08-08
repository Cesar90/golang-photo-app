package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				// return `<!-- TODO: Implement the csrfField -->`, nil
				return ``, fmt.Errorf("csrfField not implemented")
			},
		},
	)
	// tpl, err := template.ParseFS(fs, patterns...)
	tpl, err := tpl.ParseFS(fs, patterns...)

	if err != nil {
		return Template{}, fmt.Errorf("Parsing template :%w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("Parsing template :%v", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	// We clone the template to avoid race condition bugs
	// Since htmlTpl is a *template.Template (a pointer), sharing it across requests
	// can cause issues like reusing the same CSRF token across multiple request
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("Cloning template: %v, err")
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	// fmt.Fprint(w, "<h1>Welome to my awesome site!!</h1>")
	// err := t.htmlTpl.Execute(w, data)
	// err = tpl.Execute(w, data)
	err = tpl.Execute(&buf, data)
	// bytes.Buffer will stop right away if there is at least one error
	// it will give the chance to show correct error message and status code
	if err != nil {
		// panic(err) //TODO: Remove the panic
		log.Printf("parsing template %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}
