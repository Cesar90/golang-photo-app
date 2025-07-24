package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welome to my awesome site!!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at</p><a href=\"\">test</a>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ Page</h1>
		<ul>
			<li>
				Is there a free version? YES! We Offer a free trial for 30 days on any
				paid plans
			</li>
		</ul>

	`)
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
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", contactHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
