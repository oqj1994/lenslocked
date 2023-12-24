package main

import "C"
import (
	"fmt"
	"lenslocked/M"
	"lenslocked/V"
	"lenslocked/controller"
	"lenslocked/html"
	"lenslocked/migrations"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	err:=godotenv.Load(".env")
	if err != nil{
		panic(err)
	}


	db, err := M.Open(M.DefaultConfig())
	fmt.Println(M.DefaultConfig().String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = M.MigrateFS(migrations.FS, db, ".")
	if err != nil {
		panic(err)
	}


	userService := M.UserService{
		DB: db,
	}
	sessionService := M.SessionService{
		DB:            db,
		BytesPerToken: 32,
	}

	smtpPortStr:= os.Getenv("SMTP_PORT")
	smtpPort,err:=strconv.Atoi(smtpPortStr)
	if  err!=nil {
		panic(err)
	}
	emailService:=M.NewEmailService(M.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:      smtpPort,
		UserName:  os.Getenv("SMTP_USERNAME"),
		Password:  os.Getenv("SMTP_PASSWORD"),
	})

	uc := controller.UserController{
		US: userService,
		SS: sessionService,
		ES: emailService,
	}
	userMiddleware := controller.MiddleWare{SS: sessionService}
	csrfMiddleWare := csrf.Protect([]byte("abcdefghizklmnopqrstuvwxyz123456"), csrf.Secure(false))

	r := chi.NewRouter()
	r.Use(middleware.Logger, csrfMiddleWare, userMiddleware.SetUser)

	tpl := V.Must(V.ExcuteFS("index.html"))
	r.Get("/", controller.StaticController(tpl))
	tpl = V.Must(V.ExcuteFS("signup.html"))
	uc.Template.New = tpl
	tpl = V.Must(V.ExcuteFS("login.html"))
	uc.Template.Login = tpl
	r.Get("/signup", uc.New)
	r.Get("/login", uc.Login)
	r.Post("/logout", uc.Logout)
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
	r.Route("/user/me", func(r chi.Router) {

		r.Use(userMiddleware.RequireUser)
		r.Get("/", uc.CurrentUser)
	})

	fmt.Println("run server on port 10010\nPlease try to enjoy coding!!:)")
	err = http.ListenAndServe(":10010", r)
	if err != nil {
		log.Println(err)
		panic("run server error!")
	}
}
