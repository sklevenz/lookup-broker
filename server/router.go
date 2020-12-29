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
	router := mux.NewRouter()

	v2Router := router.PathPrefix("/v2").Subrouter()
	v2Router.Use(apiVersionHandler)
	v2Router.Use(requestIdentityLogHandler)
	v2Router.Use(originatingIdentityLogHandler)
	v2Router.Use(cacheHandler)
	v2Router.HandleFunc("/catalog", catalogHandler).Name("v2.catalog").Methods(http.MethodGet)
	v2Router.HandleFunc("/service_instances/{id}", instancePutHandler).Headers(headerContentType, contentTypeJSON).Name("v2.instance.put").Methods(http.MethodPut)
	v2Router.HandleFunc("/service_instances/{id}", instanceGetHandler).Name("v2.instance.get").Methods(http.MethodGet)
	v2Router.HandleFunc("/service_instances/{id}", instancePatchHandler).Headers(headerContentType, contentTypeJSON).Name("v2.instance.patch").Methods(http.MethodPatch)
	v2Router.HandleFunc("/service_instances/{id}", instanceDeleteHandler).Name("v2.instance.delete").Methods(http.MethodDelete)
	v2Router.HandleFunc("/service_instances/{id}/service_bindings/{bid}", bindingPutHandler).Headers(headerContentType, contentTypeJSON).Name("v2.binding.put").Methods(http.MethodPut)
	v2Router.HandleFunc("/service_instances/{id}/service_bindings/{bid}", bindingGetHandler).Name("v2.binding.get").Methods(http.MethodGet)
	v2Router.HandleFunc("/service_instances/{id}/service_bindings/{bid}", bindingDeleteHandler).Name("v2.binding.delete").Methods(http.MethodDelete)

	router.HandleFunc("/health", healthHandler).Name("health").Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler).Name("home").Methods(http.MethodGet)

	router.Use(logHandler)

	return router
}
