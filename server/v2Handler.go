package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sklevenz/lookup-broker/openapi"
)

const (
	supportedAPIVersionValue string = "2.14"

	headerAPIVersion            string = "X-Broker-API-Version"
	headerAPIOrginatingIdentity string = "X-Broker-API-Originating-Identity"
	headerAPIRequestIdentity    string = "X-Broker-API-Request-Identity"
)

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
	js, err := json.Marshal(buildCatalog())
	if err != nil {
		handleHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set(headerContentType, contentTypeJSON)
	reader := bytes.NewReader(js)
	http.ServeContent(w, r, "xxx", time.Now(), reader)
}

func buildCatalog() *openapi.Catalog {
	catalog := openapi.Catalog{}
	var services []openapi.Service
	var service openapi.Service

	service.Id = "lookup"
	service.Name = "lookup"
	service.Description = "Topology lookup service broker"
	service.Tags = append(service.Tags, "cf", "api", "cloudfoundry", "cloud controler", "uaa")
	service.Requires = []string{}
	service.Bindable = true
	service.InstancesRetrievable = true
	service.BindingsRetrievable = true
	service.AllowContextUpdates = true
	service.Metadata = map[string]interface{}{}
	service.DashboardClient = openapi.DashboardClient{}
	service.DashboardClient.Id = "lookupDashboardClientId"
	service.DashboardClient.RedirectUri = "https://lookup-broker.cfapps.eu10.hana.ondemand.com/"
	service.DashboardClient.Secret = "admin"
	service.PlanUpdateable = true

	plans := []openapi.Plan{}

	plan := openapi.Plan{}
	plan.Id = "cloudfoundry"
	plan.Name = "cloudfoundry"
	plan.Description = "Lookup for Cloud Controler and UAA"
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
