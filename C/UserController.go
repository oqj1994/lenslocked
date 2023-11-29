package C

import "net/http"

type UserController struct {
}

func (u UserController) Create(w http.ResponseWriter, r *http.Request) {
	//TODO: get the userName and password from request
	//TODO: use the model to add user data
	//TODO: then base on the result to render some page or return error to the responseWriter
}

func (u UserController) Find(name string) {

}
