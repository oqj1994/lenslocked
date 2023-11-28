package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
)

type IndexHandler struct {
	tpl *template.Template
}

func (h IndexHandler) Index(w http.ResponseWriter, r *http.Request) {

}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", IndexHandler{}.Index)

	fmt.Println("run server on port 10010\nPlease try to enjoy coding!!:)")
	err := http.ListenAndServe(":10010", r)
	if err != nil {
		panic("run server error!")
	}
}
