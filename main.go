package main

import (
	"fmt"
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

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h3>1.what question can I ask here?</h3>  <h4>Any dumbs question is welcome here</h4>"))
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/":
			index(writer, request)
			return
		case "/contact":
			contact(writer, request)
			return
		case "/faq":
			FaqHandler(writer, request)
		default:
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}(writer, request)
			return
		}
	})

	fmt.Println("run server on port 10010\nPlease try to enjoy coding!!:)")
	err := http.ListenAndServe(":10010", nil)
	if err != nil {
		panic("run server error!")
	}
}
