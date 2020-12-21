package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticHandler(t *testing.T) {
	log.Println("nothing to test")
}

func TestHealth(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/health/", nil)
	request.SetBasicAuth("username", "password")
	response := httptest.NewRecorder()

	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.JSONEq(t, `{"ok":true}`, response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
