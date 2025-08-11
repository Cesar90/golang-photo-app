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
	UserService    *models.UserService
	SessionService *models.SessionService
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
	u.Templates.New.Execute(w, r, data)
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

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Lon term, we should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	cookie := http.Cookie{
		Name:     "Session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// fmt.Fprint(w, "User craeted: %+v", user)
	http.Redirect(w, r, "/users/me", http.StatusFound)
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
	u.Templates.SignIn.Execute(w, r, data)
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

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true, //Cookie is not accesible from javascript
	}
	http.SetCookie(w, &cookie)
	// fmt.Fprint(w, "User authenticate: %+v", user)
	// u.Templates.SignIn.Execute(w, data)
	// u.Templates.New.Execute(w, nil)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session")
	// email, err := r.Cookie("email")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		// fmt.Fprint(w, "The email cookie could not be read.")
		return
	}
	user, err := u.SessionService.User(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		// fmt.Fprint(w, "The email cookie could not be read.")
		return
	}
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
	// fmt.Fprintf(w, "Email cookie: %s \n", email.Value)
	// fmt.Fprintf(w, "Headers: %+V \n", r.Header)
}
