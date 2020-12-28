package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	headerContentType  string = "Content-Type"
	headerETag         string = "ETag"
	headerLastModified string = "Last-Modified"

	contentTypeCSS  string = "text/css"
	contentTypeHTML string = "text/html"
	contentTypeTEXT string = "text/plain"
	contentTypeJSON string = "application/json"
)

// New implements the routes defined by OSB v2.0 API
func New() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	v2Router := router.PathPrefix("/v2/").Subrouter()
	v2Router.Use(apiVersionHandler)
	v2Router.Use(requestIdentityLogHandler)
	v2Router.Use(originatingIdentityLogHandler)
	v2Router.Use(etagHandler)
	v2Router.HandleFunc("/catalog/", catalogHandler).Name("v2.catalog").Methods(http.MethodGet)
	v2Router.HandleFunc("/service_instances/:{id}/", instancesHandler).Headers(headerContentType, contentTypeJSON).Name("v2.instances").Methods(http.MethodPut)

	router.HandleFunc("/health/", healthHandler).Name("health").Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler).Name("home").Methods(http.MethodGet)

	router.Use(logHandler)

	return router
}
