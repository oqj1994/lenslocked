package main

import "C"
import (
	"fmt"
	"lenslocked/M"
	"lenslocked/V"
	"lenslocked/controller"
	"lenslocked/html"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)



func main() {
	db, err := M.Open(M.DefaultConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	tpl := V.Must(V.ExcuteFS("index.html"))
	r.Get("/", controller.StaticController(tpl))
	tpl = V.Must(V.ExcuteFS("signin.html"))
	
	uc := controller.UserController{
		Template: struct{New controller.Template}{
			New: tpl,
		},
		US: M.UserService{
			DB: db,
		},
	}
	r.Get("/signin", uc.RenderSigninPage)
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
