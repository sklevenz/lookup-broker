package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.JSONEq(t, `{"ok":true}`, response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestHome(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeTEXT, response.Header().Get(headerContentType))
	assert.Equal(t, `Lookup-Broker`, response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
