package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func contact(w http.ResponseWriter, r *http.Request) {
	message := "call me on 131600231 :)"
	fmt.Fprint(w, message)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "This is a greate web application")
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.RawPath)
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/picture/{id}/{t}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		t := chi.URLParam(r, "t")
		fmt.Fprint(w, id)
		fmt.Fprint(w, t)
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/":
			index(writer, request)
			return
		case "/contact":
			contact(writer, request)

			return
		default:
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}(writer, request)
			return
		}
	})

	fmt.Println("run server on port 10010\nPlease try to enjoy coding!!:)")
	err := http.ListenAndServe(":10010", r)
	if err != nil {
		panic("run server error!")
	}
}
