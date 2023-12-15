package controller

import (
	"fmt"
	"lenslocked/V"
	"net/http"
)

func StaticController(tpl *V.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func ReadCookie(w http.ResponseWriter, r *http.Request) {
	cookie := r.Cookies()
	for _, c := range cookie {
		fmt.Fprintf(w, fmt.Sprintf("cookieName: %s, Value: %s \n", c.Name, c.Value))
	}
}
