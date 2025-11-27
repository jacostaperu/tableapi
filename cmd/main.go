package main

import (
	"flag"
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

	port := flag.Int("port", 8087, "port to listen on")
	tablespath := flag.String("tablespath", "tables", "where to store CSV files(remember that a valid file will contain an 'id' column)")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	help := flag.Bool("help", false, "show help")

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Use CLI flags when no config provided
	fmt.Printf("Port: %d\n", *port)
	fmt.Printf("Verbose: %v\n", *verbose)

	// Use your certificate and key files
	//certFile := "server.crt"
	//keyFile := "server.key"

	myserver := tableapi.NewServer()
	if *port < 1024 {
		myserver.Logger.Warnf("Listening to port='%d' (port < 1024) may require elevated privileges, if fail try running with sudo\n", *port)
	}

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

	myserver.SetTablesPath(*tablespath)

	myserver.Logger.Infof("Listening on port =  %d", *port)
	//return http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, handler)
	return http.ListenAndServe(fmt.Sprintf(":%d", *port), handler)

}
