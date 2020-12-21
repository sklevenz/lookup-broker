/*
 * Open Service Broker API
 *
 * The Open Service Broker API defines an HTTP(S) interface between Platforms and Service Brokers.
 *
 * API version: master - might contain changes that are not yet released
 * Contact: open-service-broker-api@googlegroups.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ServiceInstanceUpdateRequest struct {

	// See [Context Conventions](https://github.com/openservicebrokerapi/servicebroker/blob/master/profile.md#context-object) for more details.
	Context map[string]interface{} `json:"context,omitempty"`

	ServiceId string `json:"service_id"`

	PlanId string `json:"plan_id,omitempty"`

	Parameters map[string]interface{} `json:"parameters,omitempty"`

	PreviousValues ServiceInstancePreviousValues `json:"previous_values,omitempty"`

	MaintenanceInfo MaintenanceInfo `json:"maintenance_info,omitempty"`
}
