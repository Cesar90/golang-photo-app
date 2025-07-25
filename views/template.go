package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
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

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// fmt.Fprint(w, "<h1>Welome to my awesome site!!</h1>")
	err := t.htmlTpl.Execute(w, nil)
	if err != nil {
		// panic(err) //TODO: Remove the panic
		log.Printf("parsing template %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}
