package main

import "C"
import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"lenslocked/M"
	"lenslocked/V"
	"lenslocked/controller"
	"lenslocked/html"
	"log"
	"net/http"
	"path"
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
	db, err := M.Open(M.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	tpl := V.Must(V.ExcuteFS("index.html"))
	r.Get("/", controller.StaticController(tpl))
	tpl = V.Must(V.ExcuteFS("signin.html"))
	r.Get("/signin", controller.StaticController(tpl))
	uc := controller.UserController{
		US: M.UserService{
			DB: db,
		},
	}
	r.Post("/user", uc.Create)
	r.Handle("/assert/*", http.StripPrefix("/assert/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		f, err := html.FS.ReadFile(path.Join("assert", p))
		if err != nil {
			http.Error(w, "read assert error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(f)
	})))

	fmt.Println("run server on port 10010\nPlease try to enjoy coding!!:)")
	err = http.ListenAndServe(":10010", r)
	if err != nil {
		log.Println(err)
		panic("run server error!")
	}
}
