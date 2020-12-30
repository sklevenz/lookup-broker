#!/usr/bin/env bash

curl -iLs 'http://localhost:5000/v2/service_instances/:instance_id?accepts_incomplete=true' -d '{
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "service_id": "1",
  "plan_id": "1.1",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  },
  "previous_values": {
    "plan_id": "1",
    "service_id": "1.1",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}' -X PATCH -H "X-Broker-API-Version: 2.16" -H "Content-Type: application/json"
