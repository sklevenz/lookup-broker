#!/usr/bin/env bash

curl -iLs 'http://localhost:5000/v2/service_instances/123/service_bindings/456?accepts_incomplete=true' -d '{
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "service_id": "1",
  "plan_id": "1.1",
  "bind_resource": {
    "app_guid": "app-guid-here"
  },
  "parameters": {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}' -X PUT -H "X-Broker-API-Version: 2.16" -H "Content-Type: application/json"