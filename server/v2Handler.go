package server

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sklevenz/lookup-broker/landscape"
	"github.com/sklevenz/lookup-broker/openapi"
)

const (
	supportedAPIVersionValue string = "2.14"

	headerAPIVersion            string = "X-Broker-API-Version"
	headerAPIOrginatingIdentity string = "X-Broker-API-Originating-Identity"
	headerAPIRequestIdentity    string = "X-Broker-API-Request-Identity"

	catalogServiceID = "1"
	catalogPanID     = "1.1"
)

type userIDType struct {
	UserID string `json:"user_id"`
}

type originatingIdentityType struct {
	Platform string     `json:"platform"`
	UserID   userIDType `json:"user_id_object"`
}

func requestIdentityLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(headerAPIRequestIdentity)
		if headerValue != "" {
			log.Printf("Request Identity: %v", headerValue)
		} else {
			log.Printf("Header %v not set", headerAPIRequestIdentity)
		}

		next.ServeHTTP(w, r)
	})
}
func handleOSBError(w http.ResponseWriter, code int, err openapi.Error) {
	output, _ := json.Marshal(err)

	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(code)
	w.Write(output)
}
func parseOriginatingIdentityHeader(value string) (*originatingIdentityType, error) {
	value = strings.TrimSpace(value)
	values := strings.Split(value, " ")

	originatingIdentity := &originatingIdentityType{}
	originatingIdentity.Platform = values[0]

	encoded, err := base64.StdEncoding.DecodeString(values[1])

	if err != nil {
		log.Printf("Error in Originating Identity Header, user_id not base64 encoded: %v", value)
		return nil, err
	}

	json.Unmarshal([]byte(encoded), &originatingIdentity.UserID)

	return originatingIdentity, nil
}

func originatingIdentityLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(headerAPIOrginatingIdentity)
		if headerValue != "" {
			originatingIdentity, err := parseOriginatingIdentityHeader(headerValue)
			if err == nil {
				log.Printf("Originating Identity: %v", originatingIdentity)
			}
		} else {
			log.Printf("Header %v not set", headerAPIOrginatingIdentity)
		}

		next.ServeHTTP(w, r)
	})
}

func handleHTTPError(w http.ResponseWriter, code int, err error) {

	output, _ := json.Marshal(&openapi.Error{
		Error:       http.StatusText(code),
		Description: err.Error(),
	})

	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(code)
	w.Write(output)
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(buildCatalog(r))
	if err != nil {
		handleHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set(headerContentType, contentTypeJSON)
	reader := bytes.NewReader(js)
	http.ServeContent(w, r, "xxx", time.Now(), reader)
}

func buildCatalog(r *http.Request) *openapi.Catalog {
	catalog := openapi.Catalog{}
	var services []openapi.Service
	var service openapi.Service

	service.Id = catalogServiceID
	service.Name = "lookup"
	service.Description = "Lookup service broker"
	service.Tags = append(service.Tags, "cf", "api", "cloudfoundry", "cloud controler", "uaa")
	service.Requires = []string{}
	service.Bindable = true
	service.InstancesRetrievable = true
	service.BindingsRetrievable = true
	service.AllowContextUpdates = true
	service.Metadata = map[string]interface{}{}
	service.DashboardClient = openapi.DashboardClient{}
	service.DashboardClient.Id = "lookupDashboardClientId"
	service.DashboardClient.RedirectUri = ""
	service.DashboardClient.Secret = "admin"
	service.PlanUpdateable = true

	plans := []openapi.Plan{}

	plan := openapi.Plan{}
	plan.Id = catalogPanID
	plan.Name = "extension"
	plan.Description = "Topology lookup for Cloud Foundry extension landscapes"
	plan.Metadata = make(map[string]interface{})
	plan.Metadata["labels"] = []string{}
	plan.Free = true
	plan.Bindable = true
	plan.PlanUpdateable = true
	plan.Schemas = openapi.SchemasObject{}
	plan.MaximumPollingDuration = 10
	plan.MaintenanceInfo = openapi.MaintenanceInfo{}
	plan.MaintenanceInfo.Version = "0.0.0"

	plans = append(plans, plan)

	service.Plans = plans
	services = append(services, service)
	catalog.Services = services

	log.Printf("Catalog: %v", catalog)

	return &catalog
}

func cacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		landscapes := landscape.Get()
		s1 := fmt.Sprintf("%v", landscapes)
		s2 := md5.Sum([]byte(s1))
		hash := fmt.Sprintf("%x", s2)

		w.Header().Set(headerETag, fmt.Sprintf("W/\"%v\"", hash))

		next.ServeHTTP(w, r)
	})
}
func apiVersionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supportedAPIVersion := strings.Split(supportedAPIVersionValue, ".")[0]
		requestedAPIVersionValue := strings.Split(r.Header.Get(headerAPIVersion), ".")[0]

		if requestedAPIVersionValue == "" {
			err := fmt.Errorf("HTTP Status: (%v) - mandatory request header %v not set", http.StatusPreconditionFailed, headerAPIVersion)
			log.Printf("Error: %v", err)
			handleHTTPError(w, http.StatusPreconditionFailed, err)
			return
		}

		requestedAPIVersion := strings.Split(requestedAPIVersionValue, ".")[0]
		if supportedAPIVersion != requestedAPIVersion {
			err := fmt.Errorf("HTTP Status: (%v) - requested API version is %v but supported API version is %v", http.StatusPreconditionFailed, r.Header.Get(headerAPIVersion), supportedAPIVersionValue)
			log.Printf("Error: %v", err)
			handleHTTPError(w, http.StatusPreconditionFailed, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bindingPutHandler(w http.ResponseWriter, r *http.Request) {
	handleHTTPError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func bindingGetHandler(w http.ResponseWriter, r *http.Request) {
	handleHTTPError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func bindingDeleteHandler(w http.ResponseWriter, r *http.Request) {
	handleHTTPError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func instancePatchHandler(w http.ResponseWriter, r *http.Request) {
	handleHTTPError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func instanceGetHandler(w http.ResponseWriter, r *http.Request) {
	handleHTTPError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func instanceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	handleHTTPError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func instancePutHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	serviceInstanceID := vars["id"]
	log.Println("serviceInstanceId = ", serviceInstanceID)

	var requestContent openapi.ServiceInstanceProvisionRequest

	err := json.NewDecoder(r.Body).Decode(&requestContent)
	if err != nil {
		log.Printf("Error: %v", err)
		handleHTTPError(w, http.StatusBadRequest, err)
		return
	}

	if requestContent.ServiceId != catalogServiceID {
		err := errors.New("unsupported service id: " + requestContent.ServiceId)
		log.Printf("Error: %v", err)
		handleHTTPError(w, http.StatusBadRequest, err)
		return
	}

	if requestContent.PlanId != catalogPanID {
		err := errors.New("unsupported plan id: " + requestContent.PlanId)
		log.Printf("Error: %v", err)
		handleHTTPError(w, http.StatusBadRequest, err)
		return
	}

	if requestContent.OrganizationGuid == "" && requestContent.Context["organization_guid"] == "" {
		err := errors.New("organization_guid missing")
		log.Printf("Error: %v", err)
		handleHTTPError(w, http.StatusBadRequest, err)
		return
	}

	if requestContent.SpaceGuid == "" && requestContent.Context["plan_guid"] == "" {
		err := errors.New("space_guid missing")
		log.Printf("Error: %v", err)
		handleHTTPError(w, http.StatusBadRequest, err)
		return
	}

	var responseContent openapi.ServiceInstanceProvisionResponse

	js, err := json.Marshal(responseContent)
	if err != nil {
		handleHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set(headerContentType, contentTypeJSON)
	reader := bytes.NewReader(js)
	http.ServeContent(w, r, "xxx", time.Now(), reader)

	return
}
