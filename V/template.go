package V

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"lenslocked/M"
	"lenslocked/context"
	"lenslocked/html"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

func Must(t *Template, e error) *Template {
	if e != nil {
		panic(e)
	}
	return t
}

func RenderData(r *http.Request, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"data":           data,
	}
}

func ExcuteFS(name string) (*Template, error) {
	tpl := template.New("home.html")
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("this csrfField not implemented!")
		},
		"currentUser": func() error {
			return fmt.Errorf("currentUser not implement!")
		},
	})
	tpl, err := tpl.ParseFS(html.FS, "home.html", name)

	if err != nil {
		return nil, fmt.Errorf("parsing files error: %w", err)
	}
	return &Template{htmlTpl: tpl}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) error {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("clone template error: %v", err)
		return fmt.Errorf("cloning template error :%w", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl = tpl.Funcs(template.FuncMap{"csrfField": func() template.HTML {
		return csrf.TemplateField(r)
	},
		"currentUser": func() *M.User {
			user, err := context.User(r.Context())
			if err != nil {
				fmt.Println("template get user :", err)
				return nil
			}
			return user
		},
	})
	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("execute template error: %v", err)
		return fmt.Errorf("execute template error :%w", err)
	}
	io.Copy(w, &buf)
	return nil
}
