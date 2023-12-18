package controller

import (
	"fmt"
	"lenslocked/M"
	"lenslocked/V"
	"log"
	"net/http"
)

type UserController struct {
	Template struct {
		New   Template
		Login Template
	}
	US M.UserService
	SS M.SessionService
}



func (u UserController) Create(w http.ResponseWriter, r *http.Request) {
	//TODO: get the userName and password from request
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	name := "yyChat"
	user, err := u.US.Create(M.CreateUserParms{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "create user error", http.StatusInternalServerError)
		return
	}
	session, err := u.SS.Create(user.ID)
	if err != nil {
		fmt.Println("create session error ", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/user/me", http.StatusFound)
	//TODO: use the model to add user data
	//TODO: then base on the result to render some page or return error to the responseWriter
}

func (u UserController) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	user, err := u.SS.User(token)
	if err != nil {
		fmt.Println("get user error", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprint(w, user)
}

func (u UserController) Find(name string) {

}

func (u UserController) New(w http.ResponseWriter, r *http.Request) {
	err := u.Template.New.Execute(w, r, V.RenderData(r, nil))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (u UserController) Login(w http.ResponseWriter, r *http.Request) {

	err := u.Template.Login.Execute(w, r, V.RenderData(r, nil))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (u UserController) ProcessLogin(w http.ResponseWriter, r *http.Request) {
	parms := M.AuthenticateParms{}
	parms.Email = r.PostFormValue("email")
	parms.Password = r.PostFormValue("password")
	user, err := u.US.Authenticate(parms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := u.SS.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/user/me", http.StatusFound)

}
