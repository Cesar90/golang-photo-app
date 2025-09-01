package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Cesar90/golang-photo-app/context"
	"github.com/Cesar90/golang-photo-app/models"
)

type Users struct {
	// Templates struct {
	// 	New views.Template
	// }
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetSevice
	EmailService         *models.EmailService
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
	// cookie := http.Cookie{
	// 	Name:     "Session",
	// 	Value:    session.Token,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// }
	// cookie := newCookie("CookieSession", session.Token)
	// http.SetCookie(w, cookie)
	setCookie(w, "CookieSession", session.Token)
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

	// cookie := http.Cookie{
	// 	Name:     "session",
	// 	Value:    session.Token,
	// 	Path:     "/",
	// 	HttpOnly: true, //Cookie is not accesible from javascript
	// }
	// http.SetCookie(w, &cookie)
	setCookie(w, "CookieSession", session.Token)
	// fmt.Fprint(w, "User authenticate: %+v", user)
	// u.Templates.SignIn.Execute(w, data)
	// u.Templates.New.Execute(w, nil)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	// if user == nil {
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	fmt.Fprintf(w, "Current user: %\n", user.Email)

	// tokenCookie, err := r.Cookie("session")
	// token, err := readCookie(r, "CookieSession")
	// // email, err := r.Cookie("email")
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	// fmt.Fprint(w, "The email cookie could not be read.")
	// 	return
	// }
	// user, err := u.SessionService.User(token)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	// fmt.Fprint(w, "The email cookie could not be read.")
	// 	return
	// }
	// fmt.Fprintf(w, "Current user: %s\n", user.Email)
	// fmt.Fprintf(w, "Email cookie: %s \n", email.Value)
	// fmt.Fprintf(w, "Headers: %+V \n", r.Header)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, "CookieSession")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, "CookieSession")
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. For instance, if a user does not
		// exist with that email address.
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.lenslocked.com/reset-pw?" + vals.Encode()
	// "https://www.lenslocked.com/reset-pw?token=123"
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Don't render the reset token here! We need the user to confirm they have
	// access to the email account to verify their identity.
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umn UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, "CookieSession")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umn.SessionService.User(token)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
