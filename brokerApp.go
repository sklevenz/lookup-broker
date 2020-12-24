package main

import (
	"log"
	"net/http"

	"os"

	"github.com/sklevenz/lookup-broker/server"
)

const (
	defaultPort = "5000"
)

var (
	// Version set by go build via -ldflags "'-X main.Version=1.0.0'"
	Version string = "n/a"
	// Commit set by go build via -ldflags "'-X main.Commit=123'"
	Commit string = "n/a"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("start application on port: %v", port)
	log.Printf("version: %v", Version)
	log.Printf("commit: %v", Commit)

	brokerServer := server.New()

	log.Printf("call server: http://localhost:%v", port)

	if err := http.ListenAndServe(":"+port, brokerServer); err != nil {
		log.Fatalf("could not listen on port %v: %v", port, err)
	}
}
