package V

import (
	"fmt"
	"html/template"
	"lenslocked/html"
	"log"
	"net/http"
)

func Must(t *Template, e error) *Template {
	if e != nil {
		panic(e)
	}
	return t
}

func ExcuteFS(name string) (*Template, error) {
	tpl, err := template.ParseFS(html.FS, "home.html", name)
	if err != nil {
		return nil, fmt.Errorf("parsing files error: %w", err)
	}
	return &Template{htmlTpl: tpl}, nil
}

func Parse(filepath string) (*Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return nil, fmt.Errorf("parsing files error: %v", err)
	}
	return &Template{htmlTpl: tpl}, nil

}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Excute(w http.ResponseWriter, data interface{}) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("excute template error: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
