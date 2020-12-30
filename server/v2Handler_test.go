package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sklevenz/lookup-broker/openapi"
	"github.com/stretchr/testify/assert"
)

const (
	landscapes = `{
				"cf-eu10": {
				"cloudcontroller": "https://api.cf.eu10.hana.ondemand.com",
				"uaa": "https://uaa.cf.eu10.hana.ondemand.com",
				"labels": [
					"master",
					"aws"
					]
				},
				"cf-eu10-001": {
				"cloudcontroller": "https://api.cf.eu10-001.hana.ondemand.com",
				"uaa": "https://uaa.cf.eu10-001.hana.ondemand.com",
				"labels": [
					"scaleout",
					"aws"
					]
				},
				"cf-eu10-002": {
				"cloudcontroller": "https://api.cf.eu10-002.hana.ondemand.com",
				"uaa": "https://uaa.cf.eu10-002.hana.ondemand.com",
				"labels": [
					"scaleout",
					"aws"
					]
				}
			}`
)

func TestNoApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}
func TestWrongApiVersionFormat(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	request.Header.Set(headerAPIVersion, "abc")

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestWrongApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	request.Header.Set(headerAPIVersion, "1.2")

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestCorrectApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestAPIOriginatingIdentity(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	request.Header.Set(headerAPIOrginatingIdentity, "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ==")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	handler := originatingIdentityLogHandler(testHandler)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	assert.Contains(t, buf.String(), "cloudfoundry")
	assert.Contains(t, buf.String(), "683ea748-3092-4ff4-b656-39cacc4d5360")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestAPIRequestIdentity(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	request.Header.Set(headerAPIRequestIdentity, "e26cee84-6b38-4456-b34e-d1a9f002c956")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	handler := requestIdentityLogHandler(testHandler)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	assert.Contains(t, buf.String(), "e26cee84-6b38-4456-b34e-d1a9f002c956")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
func TestOSBErrorHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)

	err := &openapi.Error{
		Error:            "AsyncRequired",
		Description:      "blabla",
		InstanceUsable:   true,
		UpdateRepeatable: true,
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleOSBError(w, http.StatusInternalServerError, *err)
	})

	handler := requestIdentityLogHandler(testHandler)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Contains(t, response.Body.String(), "AsyncRequired")
	assert.Contains(t, response.Body.String(), "blabla")
	assert.Contains(t, response.Body.String(), "instance_usable")
	assert.Contains(t, response.Body.String(), "instance_usable")
}

func TestCatalogHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	request.Header.Set(headerAPIVersion, "2.2")

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Contains(t, response.Body.String(), "Lookup service broker")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, fmt.Sprintf("W/\"%v\"", "4b676504ada9e24a9f4be6c86d89a83e"), response.Header().Get(headerETag))
}

func TestInstancePutHandler(t *testing.T) {

	const payload = `{
		"service_id": "1",
		"plan_id": "1.1",
		"context": {
		  "platform": "cloudfoundry",
		  "some_field": "some-contextual-data"
		},
		"organization_guid": "org-guid-here",
		"space_guid": "space-guid-here",
		"parameters": {
		  "parameter1": 1,
		  "parameter2": "foo"
		}
	  }`

	request, err := http.NewRequest(http.MethodPut, "/v2/service_instances/123", strings.NewReader(payload))
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)
	request.Header.Set(headerContentType, contentTypeJSON)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)

	var responseContent openapi.ServiceInstanceProvisionResponse

	err = json.NewDecoder(response.Body).Decode(&responseContent)
	assert.Nil(t, err)

	assert.Equal(t, "", responseContent.DashboardUrl)
}

func TestInstancePutHandlerWrongBody(t *testing.T) {
	const payload = `{
		"service_id": "",
		"plan_id": "",
		"context": {
		  "platform": "cloudfoundry",
		  "some_field": "some-contextual-data"
		},
		"organization_guid": "org-guid-here",
		"space_guid": "space-guid-here",
		"parameters": {
		  "parameter1": 1,
		  "parameter2": "foo"
		}
	  }`

	request, err := http.NewRequest(http.MethodPut, "/v2/service_instances/123", strings.NewReader(payload))
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)
	request.Header.Set(headerContentType, contentTypeJSON)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}

func TestInstanceGetHandler(t *testing.T) {

	request, err := http.NewRequest(http.MethodGet, "/v2/service_instances/123", nil)
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	var responseContent openapi.ServiceInstanceResource

	err = json.NewDecoder(response.Body).Decode(&responseContent)
	assert.Nil(t, err)
	assert.Equal(t, "", responseContent.DashboardUrl)
}

func TestInstancePatchHandler(t *testing.T) {

	const payload = `{
		"service_id": "1",
		"plan_id": "1.1",
		"context": {
		  "platform": "cloudfoundry",
		  "some_field": "some-contextual-data"
		},
		"organization_guid": "org-guid-here",
		"space_guid": "space-guid-here",
		"parameters": {
		  "parameter1": 1,
		  "parameter2": "foo"
		}
	  }`

	request, err := http.NewRequest(http.MethodPatch, "/v2/service_instances/123", strings.NewReader(payload))
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)
	request.Header.Set(headerContentType, contentTypeJSON)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	var responseContent openapi.ServiceInstanceUpdateRequest

	err = json.NewDecoder(response.Body).Decode(&responseContent)
	assert.Nil(t, err)
}

func TestInstanceDeleteHandler(t *testing.T) {

	request, err := http.NewRequest(http.MethodDelete, "/v2/service_instances/123", nil)
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestBindingDeleteHandler(t *testing.T) {

	request, err := http.NewRequest(http.MethodDelete, "/v2/service_instances/123/service_bindings/456", nil)
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestBindingGetHandler(t *testing.T) {
	os.Setenv("LANDSCAPES", landscapes)

	request, err := http.NewRequest(http.MethodGet, "/v2/service_instances/123/service_bindings/456", nil)
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	var responseContent openapi.ServiceBindingResource

	err = json.NewDecoder(response.Body).Decode(&responseContent)

	assert.Nil(t, err)
	assert.NotNil(t, responseContent)
	assert.NotNil(t, responseContent.Parameters)
	assert.NotNil(t, responseContent.Parameters["landscapes"])
}

func TestBindingPutHandler(t *testing.T) {
	const payload = `{
		"service_id": "1",
		"plan_id": "1.1",
		"context": {
		  "platform": "cloudfoundry",
		  "some_field": "some-contextual-data"
		},
		"organization_guid": "org-guid-here",
		"space_guid": "space-guid-here",
		"parameters": {
		  "parameter1": 1,
		  "parameter2": "foo"
		}
	  }`

	os.Setenv("LANDSCAPES", landscapes)

	request, err := http.NewRequest(http.MethodPut, "/v2/service_instances/123/service_bindings/456", strings.NewReader(payload))
	assert.Nil(t, err)
	request.Header.Set(headerAPIVersion, supportedAPIVersionValue)
	request.Header.Set(headerContentType, contentTypeJSON)

	response := httptest.NewRecorder()
	New().ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)

	var responseContent openapi.ServiceBindingResponse
	err = json.NewDecoder(response.Body).Decode(&responseContent)

	assert.Nil(t, err)
	assert.NotNil(t, responseContent)
}
