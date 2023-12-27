package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"lenslocked/M"
	"lenslocked/V"
	"lenslocked/context"
	"log"
	"net/http"
	"net/url"
)

type UserController struct {
	Template struct {
		New            Template
		Login          Template
		ForgetPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	US M.UserService
	SS M.SessionService
	ES M.EmailService
	PR M.PasswordResetService
}

func (u UserController) Create(w http.ResponseWriter, r *http.Request) {
	//TODO: get the userName and password from request
	var data struct{
		Email string
		Password string
	}
	data.Email = r.PostFormValue("email")
	data.Password = r.PostFormValue("password")
	name := "yyChat"
	user, err := u.US.Create(M.CreateUserParms{
		Name:     name,
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		log.Println(err)
		u.Template.New.Execute(w,r,data,err)
		return
	}
	session, err := u.SS.Create(user.ID)
	if err != nil {
		fmt.Println("create session error ", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Println("session in controller :", session)
	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/user/me", http.StatusFound)
	//TODO: use the model to add user data
	//TODO: then base on the result to render some page or return error to the responseWriter
}

func (u UserController) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := context.User(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "current user name: %s", user.Name)
}

func (u UserController) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	user, err := u.SS.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err = u.SS.Delete(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "delete session error", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, "")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u UserController) Find(name string) {

}

func (u UserController) PasswordReset(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Template.ForgetPassword.Execute(w, r, data)
}

func (u UserController) ProcessForgetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	passwordReset, err := u.PR.Create(data.Email)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/signup", http.StatusFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {passwordReset.Token},
	}
	err = u.ES.ForgetPassword(data.Email, "http://localhost:10010/reset-pw?"+vals.Encode())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to send email", http.StatusInternalServerError)
		return
	}
	u.Template.CheckYourEmail.Execute(w, r, data)
}

// ResetPassword to render a page to reset password
func (u UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	fmt.Println("Token----------", data.Token)
	u.Template.ResetPassword.Execute(w, r, data)
}

func (u UserController) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.PostFormValue("token")
	data.Password = r.PostFormValue("password")
	user, err := u.PR.Consume(data.Token)
	if err != nil {
		//TODO: 分别讨论 token失效 还是数据库无法连接
		fmt.Println("consume error: ", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	//TODO: use UserService to update user's password
	err = u.US.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println("update user password error:", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	//TODO: if successed , login the user
	session, err := u.SS.Create(user.ID)
	if err != nil {
		fmt.Println("create session error:", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/user/me", http.StatusFound)
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

type MiddleWare struct {
	SS M.SessionService
}

func (m MiddleWare) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			fmt.Println("CurrentUser: get cookie error ", err)
			next.ServeHTTP(w, r)
			return
		}
		user, err := m.SS.User(token)
		fmt.Println(user)
		if err != nil {
			fmt.Println("get user error", err)
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m MiddleWare) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := context.User(r.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "required user", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
