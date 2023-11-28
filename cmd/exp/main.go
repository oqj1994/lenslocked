package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	tpl.Execute(w, "松塘")

}

func main() {
	r := chi.NewRouter()
	r.Get("/", indexHandler)
	fmt.Println("run on porn 9696")
	http.ListenAndServe(":9696", r)
}
