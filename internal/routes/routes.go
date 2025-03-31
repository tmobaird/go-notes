package routes

import (
	"log"
	"net/http"
)

type Route struct {
	path    string
	handler func(w http.ResponseWriter, r *http.Request)
}

func Routes() []Route {
	return []Route{
		{path: "GET /notes", handler: getNotes},
		{path: "GET /notes/new", handler: newNote},
		{path: "POST /notes", handler: createNotes},
		{path: "GET /notes/{id}", handler: showNote},
	}
}

func RegisterRoutes(mux *http.ServeMux) {
	routes := Routes()
	for _, route := range routes {
		mux.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received %s %s", r.Method, r.URL)
			route.handler(w, r)
		})
	}
}
