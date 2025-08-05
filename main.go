package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Cesar90/golang-photo-app/controllers"
	"github.com/Cesar90/golang-photo-app/models"
	"github.com/Cesar90/golang-photo-app/templates"
	"github.com/Cesar90/golang-photo-app/views"
	"github.com/go-chi/chi/v5"
)

func executionTemplate(w http.ResponseWriter, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("Parsing template %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
	}
	t.Execute(w, nil)
	// fmt.Fprint(w, "<h1>Welome to my awesome site!!</h1>")
	// tpl, err := template.ParseFiles(filepath)
	// if err != nil {
	// 	// panic(err) //TODO: Remove the panic
	// 	log.Printf("parsing template %v", err)
	// 	http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
	// 	return
	// }
	// viewTpl := views.Template{
	// 	HTMLTpl: tpl,
	// }
	// viewTpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "<h1>Welome to my awesome site!!</h1>")
	tlpPath := filepath.Join("templates", "home.gohtml")
	executionTemplate(w, tlpPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tlpPath := filepath.Join("templates", "contact.gohtml")
	executionTemplate(w, tlpPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tlpPath := filepath.Join("templates", "faq.gohtml")
	executionTemplate(w, tlpPath)
	// fmt.Fprint(w, `<h1>FAQ Page</h1>
	// 	<ul>
	// 		<li>
	// 			Is there a free version? YES! We Offer a free trial for 30 days on any
	// 			paid plans
	// 		</li>
	// 	</ul>
	// `)
}

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homeHandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	default:
// 		// w.WriteHeader(http.StatusNotFound)
// 		// fmt.Fprint(w, "Page not found")
// 		// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		http.Error(w, "Page not found", http.StatusNotFound)
// 	}
// }

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "Page not found")
		// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	// http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/contact", contactHandler)
	// http.HandleFunc("/", pathHandler)
	// var router Router
	// var router http.HandlerFunc
	// router = pathHandler
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", router)
	// http.Handler - interface with the ServeHTTP method
	// http.HandlerFunc - a function type that accepts same args as ServeHTTP method
	// also implements http.Handler type
	// http.Handle("/", http.HandlerFunc(homeHandler))
	// http.Handle("/contact", http.HandlerFunc(contactHandler))
	// http.HandleFunc("/", http.HandlerFunc(homeHandler).ServeHTTP)

	// var router Router
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", router)
	// r := chi.NewRouter()
	// r.Get("/", homeHandler)
	// r.Get("/contact", contactHandler)
	// r.Get("/faq", faqHandler)
	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "Page not found", http.StatusNotFound)
	// })
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", r)
	r := chi.NewRouter()
	// r.Get("/", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))))
	// r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "layout-parts.gohtml"))))
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml", "tailwind.gohtml",
	))))
	// r.Get("/contact", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"contact.gohtml", "tailwind.gohtml",
	))))
	// r.Get("/faq", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml", "tailwind.gohtml",
	))))

	// r.Get("/signup", controllers.FAQ(views.Must(views.ParseFS(
	// 	templates.FS,
	// 	"signup.gohtml", "tailwind.gohtml",
	// ))))
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService, //TODO: Set this!
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
