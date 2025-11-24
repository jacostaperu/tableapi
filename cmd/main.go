package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jacostaperu/tableapi.git"
	"github.com/rs/cors"
)

// Version and build time get set via ldflags
var version = "dev"
var buildTime = "unknown"
var mode = "dev"

func main() {
	log.Printf("Version: %s\nBuild Time: %s\n", version, buildTime)
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error   %s\n", err)
		os.Exit(1)
	}
}

func run() error {

	// Use your certificate and key files
	//certFile := "server.crt"
	//keyFile := "server.key"

	myserver := tableapi.NewServer()

	if mode == "prod" {
		myserver.RunProdMode()

	} else {
		myserver.RunDevMode()
	}

	mux := myserver.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://dev.bo-cc.unify.local:8089", "https://dev.bo-cc.unify.local", "http://localhost:6666"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Authorization",
			"Content-Type",
		},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(mux)
	//handler := mux
	port := 8087
	log.Printf("listening on port =  %d", port)
	//return http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, handler)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)

}
