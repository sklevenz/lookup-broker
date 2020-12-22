package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV2Handler(t *testing.T) {
	log.Println("nothing to test")
}

func TestCatalogHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	request.SetBasicAuth("username", "password")
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	New().ServeHTTP(response, request)

	assert.Contains(t, response.Body.String(), "Cloud Foundry API Service")

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	// assert.Equal(t, fmt.Sprintf("W/\"%v\"", config.GetLastModifiedHash()), response.Header().Get(headerETag))
	// assert.Equal(t, fmt.Sprintf("%v", config.GetLastModified().UTC().Format(http.TimeFormat)), response.Header().Get(headerLastModified))
}
