package main

import (
	"fmt"
	"log"
	"md-notes/internal/routes"
	"net/http"
	"os"
	"time"
)

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3000"
	}
	log.Printf("Starting server on %s", port)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
