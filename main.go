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
	"github.com/gorilla/csrf"
)

func main() {
	db, err := M.Open(M.DefaultConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	csrfMiddleWare := csrf.Protect([]byte("abcdefghizklmnopqrstuvwxyz123456"))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	tpl := V.Must(V.ExcuteFS("index.html"))
	r.Get("/", controller.StaticController(tpl))
	tpl = V.Must(V.ExcuteFS("signup.html"))

	uc := controller.UserController{
		US: M.UserService{
			DB: db,
		},
		SS: M.SessionService{
			DB:            db,
			BytesPerToken: 32,
		},
	}
	uc.Template.New = tpl
	tpl = V.Must(V.ExcuteFS("login.html"))
	uc.Template.Login = tpl
	r.Get("/signup", uc.New)
	r.Get("/login", uc.Login)
	r.Get("/user/me", uc.CurrentUser)
	r.Get("/logout", uc.Logout)
	r.Get("/cookie", controller.ReadCookie)
	r.Post("/user", uc.Create)
	r.Post("/login", uc.ProcessLogin)
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
	err = http.ListenAndServe(":10010", csrfMiddleWare(r))
	if err != nil {
		log.Println(err)
		panic("run server error!")
	}
}
