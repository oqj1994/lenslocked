package controller

import "net/http"

const CookieSession = "session"

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
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
