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
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
		// CSRFField template.HTML
	}
	//Read variables from url
	//http://localhost:3000/signup?email=test@test.test
	email := r.FormValue("email")
	if email != "" {
		data.Email = email
	}
	// data.CSRFField = csrf.TemplateField(r)
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
	// fmt.Fprint(w, "Email:", r.FormValue("email"))
	// fmt.Fprint(w, "Password:", r.FormValue("password"))
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "User craeted: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
		// CSRFField template.HTML
	}
	//Read variables from url
	//http://localhost:3000/signup?email=test@test.test
	email := r.FormValue("email")
	if email != "" {
		data.Email = email
	}
	// data.CSRFField = csrf.TemplateField(r)
	u.Templates.SignIn.Execute(w, data)
	// u.Templates.New.Execute(w, nil)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true, //Cookie is not accesible from javascript
	}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w, "User authenticate: %+v", user)
	// u.Templates.SignIn.Execute(w, data)
	// u.Templates.New.Execute(w, nil)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "The email cookie could not be read.")
		return
	}
	fmt.Fprintf(w, "Email cookie: %s \n", email.Value)
	fmt.Fprintf(w, "Headers: %+V \n", r.Header)
}
