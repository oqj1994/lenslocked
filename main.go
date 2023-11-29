package main

import "C"
import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"lenslocked/V"
	"lenslocked/controller"
	"log"
	"net/http"
	"path/filepath"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	t, err := V.Parse(filepath)
	if err != nil {
		log.Printf("parse files error %v", err)
		http.Error(w, "parse files error", http.StatusInternalServerError)
		return
	}
	t.Excute(w, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("html", "index.html")
	executeTemplate(w, path)

}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	tpl, err := V.Parse(filepath.Join("html", "index.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controller.StaticController(tpl))

	fmt.Println("run server on port 10010\nPlease try to enjoy coding!!:)")
	err = http.ListenAndServe(":10010", r)
	if err != nil {
		panic("run server error!")
	}
}
