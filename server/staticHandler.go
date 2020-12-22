package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeJSON)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeTEXT)
	w.Write([]byte("Lookup-Broker"))
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- new request ------------------------------------")
		log.Printf("raw request object: %v", r)

		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("execution time: %v", time.Since(start))
	})
}
