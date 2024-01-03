package controller

import (
	"fmt"
	"lenslocked/M"
	"lenslocked/context"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GalleryController struct {
	Templates struct {
		New  Template
		Home Template
		Edit Template
		List Template
	}
	GS M.GalleryService
}

func (gc GalleryController) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID int
	}
	user, err := context.User(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	data.ID = user.ID
	gc.Templates.New.Execute(w, r, data, nil)
}

func (gc GalleryController) Create(w http.ResponseWriter, r *http.Request) {
	user, err := context.User(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	title := r.PostFormValue("title")
	desciption := r.PostFormValue("desciption")
	_, err = gc.GS.Create(title, desciption, user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/gallery/home", http.StatusFound)
}

func (gc GalleryController) Home(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Gallerys []M.Gallery
	}
	user, err := context.User(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	data.Gallerys, err = gc.GS.List(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	gc.Templates.Home.Execute(w, r, data)
}

func (gc GalleryController) Delete(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	galleryID, err := strconv.Atoi(IDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "url params error ", http.StatusBadRequest)
		return
	}
	err = gc.GS.Delete(galleryID)
	if err != nil {
		log.Println(err)
		http.Error(w, "delete gallery error ", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/gallery/home", http.StatusFound)
}

func (gc GalleryController) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := context.Gallery(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/gallery/home", http.StatusFound)
		return
	}
	gc.Templates.Edit.Execute(w, r, gallery)
}

func (gc GalleryController) List(w http.ResponseWriter, r *http.Request) {

	type Image struct {
		FileName  string
		Path      string
		GalleryID int
	}
	IDStr := chi.URLParam(r, "id")
	galleryID, err := strconv.Atoi(IDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "url params error ", http.StatusBadRequest)
		return
	}
	gallery, err := gc.GS.ByID(galleryID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/gallery/home", http.StatusFound)
		return
	}
	var data struct {
		Gallery M.Gallery
		Images  []Image
	}
	data.Gallery.Title = gallery.Title
	data.Gallery.Desciption = gallery.Desciption
	imgs, err := gc.GS.Images(galleryID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	for _, img := range imgs {
		data.Images = append(data.Images, Image{FileName: img.FileName, GalleryID: galleryID})
	}

	gc.Templates.List.Execute(w, r, data)
}

func (gc GalleryController) Image(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	galleryID, err := strconv.Atoi(IDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "url params error ", http.StatusBadRequest)
		return
	}
	fileName := chi.URLParam(r, "filename")
	dir := filepath.Join(gc.GS.GalleryDir(galleryID), fileName)
	http.ServeFile(w, r, dir)
}

func (gc GalleryController) Update(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	galleryID, err := strconv.Atoi(IDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "url params error ", http.StatusBadRequest)
		return
	}
	title := r.PostFormValue("title")
	description := r.PostFormValue("desciption")
	err = gc.GS.Update(title, description, galleryID)
	if err != nil {
		log.Println(err)
		http.Error(w, "gallery update failed ", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/gallery/home", http.StatusFound)
}

type GalleryMiddleware struct {
	GS M.GalleryService
}

func (gm GalleryMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get user from context
		user, err := context.User(r.Context())
		if err != nil {
			log.Println(err)
			http.Error(w, "permission denide", http.StatusBadRequest)
			return
		}
		//get gallery id from url param
		IDStr := chi.URLParam(r, "id")
		galleryID, err := strconv.Atoi(IDStr)
		if err != nil {
			log.Println(err)
			http.Error(w, "url params error ", http.StatusBadRequest)
			return
		}
		//auth the user has the gallery
		if gm.GS.UserHave(user.ID, galleryID) != nil {
			log.Printf("user(%d) don't have permission to access gallery(%d)", user.ID, galleryID)
			http.Error(w, "permission denide", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (gm GalleryMiddleware) GalleryRequire(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		IDStr := chi.URLParam(r, "id")
		galleryID, err := strconv.Atoi(IDStr)
		if err != nil {
			log.Println(err)
			http.Error(w, "url params error ", http.StatusBadRequest)
			return
		}
		gallery, err := gm.GS.ByID(galleryID)
		if err != nil {
			log.Println(err)
			http.Error(w, "get gallery error ", http.StatusInternalServerError)
			return
		}
		ctx := context.WithGallery(r.Context(), gallery)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
