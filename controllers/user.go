package controllers

import (
	"fmt"
	"net/http"

	"github.com/Cesar90/golang-photo-app/models"
)

type Users struct {
	// Templates struct {
	// 	New views.Template
	// }
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	//Read variables from url
	//http://localhost:3000/signup?email=test@test.test
	email := r.FormValue("email")
	if email != "" {
		data.Email = email
	}
	u.Templates.New.Execute(w, data)
	// u.Templates.New.Execute(w, nil)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "Tempory response")
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// fmt.Fprint(w, "Email:", r.PostForm.Get("email"))
	// fmt.Fprint(w, "Password:", r.PostForm.Get("password"))
	fmt.Fprint(w, "Email:", r.FormValue("email"))
	fmt.Fprint(w, "Password:", r.FormValue("password"))

}
