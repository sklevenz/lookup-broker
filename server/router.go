package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	headerContentType  string = "Content-Type"
	headerETag         string = "ETag"
	headerLastModified string = "Last-Modified"

	contentTypeCSS  string = "text/css; charset=utf-8"
	contentTypeHTML string = "text/html; charset=utf-8"
	contentTypeTEXT string = "text/plain; charset=utf-8"
	contentTypeJSON string = "application/json; charset=utf-8"
)

// New implements the routes defined by OSB v2.0 API
func New() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	v2Router := router.PathPrefix("/v2/").Subrouter()
	v2Router.HandleFunc("/catalog/", catalogHandler).Name("v2.catalog").Methods(http.MethodGet)

	router.HandleFunc("/health/", healthHandler).Name("health").Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler).Name("home").Methods(http.MethodGet)

	router.Use(logHandler)

	return router
}
