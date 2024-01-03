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
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type config struct {
	SMTPConfig     M.SMTPConfig
	PostgresConfig M.PostgresConfig
	Server         struct {
		Address string
		Port    int
	}
	CSRF struct {
		Key    string
		Secure bool
	}
}

func initConfig() (config, error) {
	var cfg config
	err := godotenv.Load(".env")
	if err != nil {
		return cfg, err
	}
	cfg.PostgresConfig = M.DefaultConfig()
	cfg.SMTPConfig.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTPConfig.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTPConfig.UserName = os.Getenv("SMTP_USERNAME")
	cfg.SMTPConfig.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "abcdefghizklmnopqrstuvwxyz123456"
	cfg.CSRF.Secure = false

	serverPortStr := os.Getenv("SERVER_PORT")
	cfg.Server.Port, err = strconv.Atoi(serverPortStr)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		panic(err)
	}

	db, err := M.Open(cfg.PostgresConfig)
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
	passwordResetService := M.PasswordResetService{
		BytesPerToken: 32,
		DB:            db,
		Duration:      1 * time.Hour,
	}

	emailService := M.NewEmailService(cfg.SMTPConfig)

	uc := controller.UserController{
		US: userService,
		SS: sessionService,
		ES: emailService,
		PR: passwordResetService,
	}
	userMiddleware := controller.MiddleWare{SS: sessionService}
	csrfMiddleWare := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

	r := chi.NewRouter()
	r.Use(middleware.Logger, csrfMiddleWare, userMiddleware.SetUser)

	r.Get("/", controller.StaticController(V.Must(V.ExcuteFS("index.html"))))
	uc.Template.New = V.Must(V.ExcuteFS("signup.html"))
	uc.Template.Login = V.Must(V.ExcuteFS("login.html"))
	uc.Template.ForgetPassword = V.Must(V.ExcuteFS("forgetpassword.html"))
	uc.Template.CheckYourEmail = V.Must(V.ExcuteFS("checkemail.html"))
	uc.Template.ResetPassword = V.Must(V.ExcuteFS("resetpassword.html"))

	//GalleryService init

	gs := M.GalleryService{DB: db}
	galleryMiddle := controller.GalleryMiddleware{
		GS: gs,
	}
	gc := controller.GalleryController{
		GS: gs,
	}
	gc.Templates.New = V.Must(V.ExcuteFS("newGallery.html"))
	gc.Templates.Home = V.Must(V.ExcuteFS("galleryHome.html"))
	gc.Templates.Edit = V.Must(V.ExcuteFS("galleryEdit.html"))
	gc.Templates.List = V.Must(V.ExcuteFS("galleryIndex.html"))

	r.Route("/gallery", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(userMiddleware.RequireUser)
			r.Get("/new", gc.New)
			r.Post("/new", gc.Create)
			r.Get("/home", gc.Home)
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/images/{filename}", gc.Image)
			r.Group(func(r chi.Router) {
				r.Use(galleryMiddle.GalleryRequire)
				r.Use(galleryMiddle.Auth)

				r.Post("/update", gc.Update)
				r.Get("/edit", gc.Edit)
				r.Post("/delete", gc.Delete)
			})
			r.Get("/", gc.List)

		})

	})

	r.Get("/signup", uc.New)
	r.Get("/login", uc.Login)
	r.Get("/forgetPW", uc.PasswordReset)
	r.Get("/reset-pw", uc.ResetPassword)
	r.Post("/process-reset-pw", uc.ProcessResetPassword)
	r.Post("/precessForgetPassword", uc.ProcessForgetPassword)
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

	fmt.Printf("run server on port %d\nPlease try to enjoy coding!!:)", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), r)
	if err != nil {
		log.Println(err)
		panic("run server error!")
	}
}
