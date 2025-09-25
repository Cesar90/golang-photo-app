package controllers

import (
	"net/http"

	"github.com/Cesar90/golang-photo-app/models"
)

type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	//Read title variable from url
	//http://localhost:3000/signup?title=Many%20Many%20cats
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}
