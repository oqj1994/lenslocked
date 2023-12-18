package controller

import (
	"net/http"
	"time"
)

const (
	CookieSession         = "session"
	DefaultCookieLiveTime = 60 * 60
)

func newCookie(name, value string) *http.Cookie {
	expiredTime := time.Now().Add(DefaultCookieLiveTime * time.Second)
	if value == "" {
		expiredTime = time.Now()
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiredTime,
	}
}

func setCookie(w http.ResponseWriter, name, value string) {
	cookie := newCookie(name, value)
	http.SetCookie(w, cookie)
}

func readCookie(r *http.Request, name string) (string, error) {
	TokenCookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return TokenCookie.Value, nil
}
