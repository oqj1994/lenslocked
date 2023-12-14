package controller

import (
	"lenslocked/V"
	"net/http"
)

func StaticController(tpl *V.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}
