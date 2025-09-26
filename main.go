package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Cesar90/golang-photo-app/controllers"
	"github.com/Cesar90/golang-photo-app/migrations"
	"github.com/Cesar90/golang-photo-app/models"
	"github.com/Cesar90/golang-photo-app/templates"
	"github.com/Cesar90/golang-photo-app/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

func executionTemplate(w http.ResponseWriter, r *http.Request, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("Parsing template %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
	}
	t.Execute(w, r, nil)
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
	executionTemplate(w, r, tlpPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tlpPath := filepath.Join("templates", "contact.gohtml")
	executionTemplate(w, r, tlpPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tlpPath := filepath.Join("templates", "faq.gohtml")
	executionTemplate(w, r, tlpPath)
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

type config struct {
	PSQL models.PostgresConfig
	SMTP models.STMPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	// TODO: Read the PSQL values from an ENV variables
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: Read the CSRF values from an ENV variable
	// TODO: Read the server values from an ENV variable
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnBeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":3000"
	return cfg, nil
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
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	///////////////////////////////
	// Setup the database
	// cfg := models.DefaultPostgresConfig()
	// fmt.Println(cfg.String())
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// err = models.Migrate(db, "migrations")
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	pwResetService := models.PasswordResetSevice{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	galleryService := &models.GalleryService{
		DB: db,
	}

	// Setup middleware
	umn := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	// crsfKey := "gFvi45R4fy5xNBlnBeZtQbfAVCYEIAUX"
	// csrfKey := []byte("gFvi45R4fy5xNBlnBeZtQbfAVCYEIAUX")
	csrfMv := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		// csrf.TrustedOrigins([]string{"http://localhost:3000", "http://127.0.0.1:3000"}),
		//This means the CSRF cookie is valid for the entire site
		csrf.Path("/"),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:          &userService, //TODO: Set this!
		SessionService:       sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml",
	))

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.gohtml", "tailwind.gohtml",
	))

	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries/edit.gohtml", "tailwind.gohtml",
	))

	r := chi.NewRouter()
	r.Use(csrfMv)
	r.Use(umn.SetUser)
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

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	// r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umn.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	// r.Get("/galleries/new", galleriesC.New)
	r.Route("/galleries", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// This middleware will apply rules to set of
			// routes
			r.Use(umn.RequireUser)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	// fmt.Println("Starting the server on :3000...")
	fmt.Printf("Starting the server on %s... \n", cfg.Server.Address)
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", csrfMv(umn.SetUser(r)))
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
	// http.ListenAndServe(":3000", TimerMiddleware(r.ServeHTTP))
}

func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		fmt.Println("Request time:", time.Since(start))
	}
}
