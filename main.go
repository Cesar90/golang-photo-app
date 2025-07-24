package main

import (
	"fmt"
	"net/http"
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
	var router Router
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
